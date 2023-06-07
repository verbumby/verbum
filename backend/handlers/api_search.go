package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/verbumby/verbum/backend/article"
	"github.com/verbumby/verbum/backend/chttp"
	"github.com/verbumby/verbum/backend/dictionary"
	"github.com/verbumby/verbum/backend/htmlui"
	"github.com/verbumby/verbum/backend/storage"
)

// APISearch search endpoint
func APISearch(w http.ResponseWriter, rctx *chttp.Context) error {
	urlQuery := htmlui.Query([]htmlui.QueryParam{
		htmlui.NewStringQueryParam("q", ""),
		htmlui.NewInDictsQueryParam("in"),
		htmlui.NewIntegerQueryParam("page", 1),
	})
	urlQuery.From(rctx.R.URL.Query())

	q := urlQuery.Get("q").(*htmlui.StringQueryParam).Value()
	if len(q) > 1000 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil
	}
	page := urlQuery.Get("page").(*htmlui.IntegerQueryParam).Value()
	inDicts := urlQuery.Get("in").(*htmlui.InDictsQueryParam).Value()

	inDictsStr := ""
	for _, d := range inDicts {
		if len(inDictsStr) == 0 {
			inDictsStr = "dict-" + dictionary.Get(d).IndexID()
		} else {
			inDictsStr += ",dict-" + dictionary.Get(d).IndexID()
		}
	}

	const pageSize = 10
	reqbody := map[string]interface{}{
		"from": (page - 1) * pageSize,
		"size": pageSize,
		"query": map[string]interface{}{
			"simple_query_string": map[string]any{
				"query": q,
				"fields": []string{
					"Headword^5",
					"Headword.Smaller^4",
					"HeadwordAlt^3",
					"HeadwordAlt.Smaller^2",
					"Phrases^1",
					"Content^0",
				},
				"default_operator": "AND",
			},
		},
		"highlight": map[string]any{
			"fields": map[string]any{
				"Content": map[string]any{},
			},
			"number_of_fragments": 0,

			"pre_tags":  []string{"<highlight>"},
			"post_tags": []string{"</highlight>"},
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
				Source    article.Article `json:"_source"`
				Index     string          `json:"_index"`
				ID        string          `json:"_id"`
				Highlight struct {
					Content []string
				} `json:"highlight"`
			} `json:"hits"`
		} `json:"hits"`
		Suggest map[string][]struct {
			Options []struct {
				Text string `json:"text"`
			} `json:"options"`
		} `json:"suggest"`
	}{}

	if err := storage.Post("/"+inDictsStr+"/_search", reqbody, &respbody); err != nil {
		return fmt.Errorf("query elastic: %w", err)
	}

	type articleview struct {
		ID           string
		Content      string
		DictionaryID string
	}

	articleviews := []articleview{}
	for _, hit := range respbody.Hits.Hits {
		indexID := strings.TrimPrefix(hit.Index, "dict-")
		dict := dictionary.GetByIndexID(indexID)

		content := hit.Source.Content

		if len(hit.Highlight.Content) == 0 {
			log.Printf("APISearch: no highlights for %s/%s queried with `%s` in %s ", indexID, hit.ID, q, inDictsStr)
		} else if len(hit.Highlight.Content) > 1 {
			log.Printf("APISearch: more than 1 highlight for %s/%s queried with %s in %s ", indexID, hit.ID, q, inDictsStr)
		} else {
			content = hit.Highlight.Content[0]
			content = strings.ReplaceAll(content, "[']<highlight>", "<highlight>[']")
		}

		content = string(dict.ToHTML(content))

		articleviews = append(articleviews, articleview{
			ID:           hit.ID,
			Content:      content,
			DictionaryID: dict.ID(),
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
