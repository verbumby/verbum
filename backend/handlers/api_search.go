package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/verbumby/verbum/backend/article"
	"github.com/verbumby/verbum/backend/chttp"
	"github.com/verbumby/verbum/backend/dictionary"
	"github.com/verbumby/verbum/backend/htmlui"
	"github.com/verbumby/verbum/backend/storage"
	"github.com/verbumby/verbum/backend/textutil"
)

var searchFields = []string{
	"Headword^5",
	"Headword.Smaller^4",
	"HeadwordAlt^3",
	"HeadwordAlt.Smaller^2",
	"Phrases^1",
	"Content^1",
}

// APISearch search endpoint
func APISearch(w http.ResponseWriter, rctx *chttp.Context) error {
	urlQuery := htmlui.Query([]htmlui.QueryParam{
		htmlui.NewStringQueryParam("q", ""),
		htmlui.NewInDictsQueryParam("in"),
		htmlui.NewIntegerQueryParam("page", 1),
		htmlui.NewStringQueryParam("prefix", ""),
		htmlui.NewBoolQueryParam("track_total_hits", false),
	})
	urlQuery.From(rctx.R.URL.Query())

	queryBoolMusts := []map[string]any{}

	q := urlQuery.Get("q").(*htmlui.StringQueryParam).Value()
	if len(q) > 1000 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil
	}
	q = textutil.NormalizeQuery(q)

	if q != "" {
		queryBoolMusts = append(queryBoolMusts, map[string]any{
			"bool": map[string]any{
				"should": []map[string]any{
					{
						"match": map[string]any{
							"Headword": map[string]any{
								"query": q,
								"boost": 2.0,
							},
						},
					},
					{
						"simple_query_string": map[string]any{
							"query":            q,
							"fields":           searchFields,
							"default_operator": "AND",
						},
					},
				},
			},
		})
	}

	page := urlQuery.Get("page").(*htmlui.IntegerQueryParam).Value()
	inDicts := urlQuery.Get("in").(*htmlui.InDictsQueryParam).Value()

	inDictsStr := ""
	indicesBoost := []map[string]float32{}
	for _, dn := range inDicts {
		d := dictionary.Get(dn)
		indexName := "dict-" + d.IndexID()
		if len(inDictsStr) == 0 {
			inDictsStr = indexName
		} else {
			inDictsStr += "," + indexName
		}
		indicesBoost = append(indicesBoost, map[string]float32{
			indexName: d.Boost(),
		})
	}

	prefix := []rune(urlQuery.Get("prefix").(*htmlui.StringQueryParam).Value())
	prefix = prefix[:min(3, len(prefix))]

	if len(prefix) > 0 {
		prefixMusts := []any{}
		for i, p := range prefix {
			prefixMusts = append(prefixMusts, map[string]any{
				"term": map[string]any{
					fmt.Sprintf("Prefix.Letter%d", i+1): string(p),
				},
			})
		}

		queryBoolMusts = append(queryBoolMusts, map[string]any{
			"nested": map[string]any{
				"path": "Prefix",
				"query": map[string]any{
					"bool": map[string]any{
						"must": prefixMusts,
					},
				},
			},
		})
	}

	const pageSize = 10
	reqbody := map[string]any{
		"from": (page - 1) * pageSize,
		"size": pageSize,
		"query": map[string]any{
			"bool": map[string]any{
				"must": queryBoolMusts,
			},
		},
		"indices_boost": indicesBoost,
		"highlight": map[string]any{
			"fields": map[string]any{
				"Content": map[string]any{},
			},
			"number_of_fragments": 0,

			"pre_tags":  []string{"<highlight>"},
			"post_tags": []string{"</highlight>"},
		},
		"suggest": map[string]any{
			"text": q,
			"OverHeadword": map[string]any{
				"term": map[string]any{
					"field":         "Headword",
					"size":          5,
					"prefix_length": 0,
				},
			},
			"OverHeadwordSmaller": map[string]any{
				"term": map[string]any{
					"field":         "Headword.Smaller",
					"size":          5,
					"prefix_length": 0,
				},
			},
			"OverHeadwordAltSmaller": map[string]any{
				"term": map[string]any{
					"field":         "HeadwordAlt.Smaller",
					"size":          5,
					"prefix_length": 0,
				},
			},
			"OverHeadwordAlt": map[string]any{
				"term": map[string]any{
					"field":         "HeadwordAlt",
					"size":          5,
					"prefix_length": 0,
				},
			},
		},
	}

	if q == "" {
		reqbody["sort"] = []any{"SortKey"}
	}

	if urlQuery.Get("track_total_hits").(*htmlui.BoolQueryParam).Value() {
		reqbody["track_total_hits"] = true
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
		dict := dictionary.GetByIndex(hit.Index)
		if dict == nil {
			return fmt.Errorf("can't find dict by index %s", hit.Index)
		}

		content := hit.Source.Content

		if len(hit.Highlight.Content) == 1 {
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
		DictIDs         []string
		Q               string
		Prefix          string
		Articles        []articleview
		TermSuggestions []string
		Pagination      paginationview
	}{
		DictIDs:         inDicts,
		Q:               q,
		Prefix:          string(prefix),
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
