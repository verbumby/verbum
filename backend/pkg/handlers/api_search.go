package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/verbumby/verbum/backend/pkg/article"
	"github.com/verbumby/verbum/backend/pkg/chttp"
	"github.com/verbumby/verbum/backend/pkg/dictionary"
	"github.com/verbumby/verbum/backend/pkg/htmlui"
	"github.com/verbumby/verbum/backend/pkg/storage"
)

// APISearch search endpoint
func APISearch(w http.ResponseWriter, rctx *chttp.Context) error {
	urlQuery := htmlui.Query([]htmlui.QueryParam{
		htmlui.NewStringQueryParam("q", ""),
		htmlui.NewIntegerQueryParam("page", 1),
	})
	urlQuery.From(rctx.R.URL.Query())

	q := urlQuery.Get("q").(*htmlui.StringQueryParam).Value()
	if len(q) > 1000 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil
	}
	page := urlQuery.Get("page").(*htmlui.IntegerQueryParam).Value()

	const pageSize = 10
	reqbody := map[string]interface{}{
		"from": (page - 1) * pageSize,
		"size": pageSize,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": q,
				"fields": []string{
					"Headword^4",
					"Headword.Smaller^3",
					"HeadwordAlt^2",
					"HeadwordAlt.Smaller^1",
				},
			},
		},
		"suggest": map[string]interface{}{
			"text": q,
			"OverHeadword": map[string]interface{}{
				"term": map[string]interface{}{
					"field":         "Headword",
					"size":          5,
					"prefix_length": 0,
				},
			},
			"OverHeadwordSmaller": map[string]interface{}{
				"term": map[string]interface{}{
					"field":         "Headword.Smaller",
					"size":          5,
					"prefix_length": 0,
				},
			},
			"OverHeadwordAltSmaller": map[string]interface{}{
				"term": map[string]interface{}{
					"field":         "HeadwordAlt.Smaller",
					"size":          5,
					"prefix_length": 0,
				},
			},
			"OverHeadwordAlt": map[string]interface{}{
				"term": map[string]interface{}{
					"field":         "HeadwordAlt",
					"size":          5,
					"prefix_length": 0,
				},
			},
		},
	}

	respbody := struct {
		Hits struct {
			Total struct {
				Value    int    `json:"value"`
				Relation string `json:"relation"`
			} `json:"total"`
			Hits []struct {
				Source article.Article `json:"_source"`
				Index  string          `json:"_index"`
				ID     string          `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
		Suggest map[string][]struct {
			Options []struct {
				Text string `json:"text"`
			} `json:"options"`
		} `json:"suggest"`
	}{}

	if err := storage.Post("/dict-*/_search", reqbody, &respbody); err != nil {
		return fmt.Errorf("query elastic: %w", err)
	}

	articles := []article.Article{}
	dicts := dictionary.GetAllAsMap()
	for _, hit := range respbody.Hits.Hits {
		dictID := strings.TrimPrefix(hit.Index, "dict-")

		article := hit.Source
		article.ID = hit.ID
		article.Dictionary = dicts[dictID]
		articles = append(articles, article)
	}

	type articleview struct {
		ID           string
		Content      string
		DictionaryID string
	}

	articleviews := []articleview{}
	for _, a := range articles {
		articleviews = append(articleviews, articleview{
			ID:           a.ID,
			Content:      string(a.Dictionary.ToHTML(a.Content, a.Title)),
			DictionaryID: a.Dictionary.ID(),
		})
	}

	termSuggestions := []string{}
	termSuggestionsMap := map[string]bool{}
	if len(articleviews) == 0 {
	outer:
		for _, suggestsByField := range respbody.Suggest {
			for _, suggestByField := range suggestsByField {
				for _, option := range suggestByField.Options {
					if _, ok := termSuggestionsMap[option.Text]; !ok {
						termSuggestionsMap[option.Text] = true
						termSuggestions = append(termSuggestions, option.Text)
						if len(termSuggestions) == 5 {
							break outer
						}
					}
				}
			}
		}
	}

	type paginationview struct {
		Current  int
		Total    int
		Relation string
	}

	if err := json.NewEncoder(w).Encode(struct {
		Articles        []articleview
		TermSuggestions []string
		Pagination      paginationview
	}{
		Articles:        articleviews,
		TermSuggestions: termSuggestions,
		Pagination: paginationview{
			Current:  urlQuery.Get("page").(*htmlui.IntegerQueryParam).Value(),
			Total:    int(math.Ceil(float64(respbody.Hits.Total.Value) / pageSize)),
			Relation: respbody.Hits.Total.Relation,
		},
	}); err != nil {
		return fmt.Errorf("encode response: %w", err)
	}

	return nil
}
