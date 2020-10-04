package handlers

import (
	"fmt"
	"math"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/verbumby/verbum/backend/pkg/article"
	"github.com/verbumby/verbum/backend/pkg/chttp"
	"github.com/verbumby/verbum/backend/pkg/dictionary"
	"github.com/verbumby/verbum/backend/pkg/htmlui"
	"github.com/verbumby/verbum/backend/pkg/storage"
	"github.com/verbumby/verbum/backend/pkg/tm"
)

// Dictionary handles dictionary page request
func Dictionary(w http.ResponseWriter, rctx *chttp.Context) error {
	vars := mux.Vars(rctx.R)
	dictID := vars["dictionary"]

	dict := dictionary.Get(dictID)
	if dict == nil {
		return fmt.Errorf("dictionary get %s: not found", dictID)
	}

	urlQuery := htmlui.Query([]htmlui.QueryParam{
		htmlui.NewIntegerQueryParam("page", 1),
		htmlui.NewStringQueryParam("prefix", ""),
	})
	urlQuery.From(rctx.R.URL.Query())

	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	prefix := []rune(urlQuery.Get("prefix").(*htmlui.StringQueryParam).Value())
	prefix = prefix[:min(3, len(prefix))]
	allaggs := map[string]interface{}{}
	for i := 0; i < min(len(prefix), 2)+1; i++ {
		aggs := map[string]interface{}{
			fmt.Sprintf("Letter%d", i+1): map[string]interface{}{
				"terms": map[string]interface{}{
					"field": fmt.Sprintf("Prefix.Letter%d", i+1),
					"size":  200,
					"order": map[string]interface{}{"_key": "asc"},
				},
			},
		}
		for j := 0; j < i; j++ {
			aggs = map[string]interface{}{
				fmt.Sprintf("Letter%d", i+1): map[string]interface{}{
					"filter": map[string]interface{}{
						"term": map[string]interface{}{
							fmt.Sprintf("Prefix.Letter%d", j+1): string(prefix[j]),
						},
					},
					"aggs": aggs,
				},
			}
		}

		for k, v := range aggs {
			allaggs[k] = v
		}
	}

	aggsreqbody := map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			"Prefix": map[string]interface{}{
				"nested": map[string]interface{}{"path": "Prefix"},
				"aggs":   allaggs,
			},
		},
	}

	type aggresult struct {
		Buckets []htmlui.LetterFilterEntity `json:"buckets"`
	}
	aggsrespbody := struct {
		Aggregations struct {
			Prefix struct {
				Letter1 *aggresult
				Letter2 *struct{ Letter2 aggresult }
				Letter3 *struct{ Letter3 struct{ Letter3 aggresult } }
			}
		} `json:"aggregations"`
	}{}
	if err := storage.Post("/dict-"+dictID+"/_search", aggsreqbody, &aggsrespbody); err != nil {
		return fmt.Errorf("aggs query: %w", err)
	}
	letterFilter := htmlui.LetterFilter{
		Prefix: prefix,
		LetterLink: func(prefix string) string {
			urlQuery := urlQuery.Clone()
			urlQuery.Get("page").Reset()
			urlQuery.Get("prefix").(*htmlui.StringQueryParam).SetValue(prefix)
			result := urlQuery.Encode()
			if result != "" {
				result = "?" + result
			}
			return rctx.R.URL.Path + result
		},
	}
	letterFilter.AddLevel(aggsrespbody.Aggregations.Prefix.Letter1.Buckets)
	if aggsrespbody.Aggregations.Prefix.Letter2 != nil {
		letterFilter.AddLevel(aggsrespbody.Aggregations.Prefix.Letter2.Letter2.Buckets)
	}
	if aggsrespbody.Aggregations.Prefix.Letter3 != nil {
		letterFilter.AddLevel(aggsrespbody.Aggregations.Prefix.Letter3.Letter3.Letter3.Buckets)
	}

	const pageSize = 10
	prefixMusts := []interface{}{}
	for i, p := range prefix {
		prefixMusts = append(prefixMusts, map[string]interface{}{
			"term": map[string]interface{}{
				fmt.Sprintf("Prefix.Letter%d", i+1): string(p),
			},
		})
	}
	articles, total, err := article.Query("/dict-"+dictID+"/_search", map[string]interface{}{
		"from": (urlQuery.Get("page").(*htmlui.IntegerQueryParam).Value() - 1) * pageSize,
		"size": pageSize,
		"sort": []interface{}{
			"Title",
		},
		"query": map[string]interface{}{
			"nested": map[string]interface{}{
				"path": "Prefix",
				"query": map[string]interface{}{
					"bool": map[string]interface{}{
						"must": prefixMusts,
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("query articles: %w", err)
	}

	err = tm.Render("dictionary", w, struct {
		PageTitle       string
		PageDescription string
		MetaRobotsTag   htmlui.MetaRobotsTag
		Dictionary      dictionary.Dictionary
		Articles        []article.Article
		Pagination      htmlui.Pagination
		LetterFilter    htmlui.LetterFilter
	}{
		PageTitle:       dict.Title(),
		PageDescription: dict.Title(),
		MetaRobotsTag:   htmlui.MetaRobotsTag{Index: false, Follow: true},
		Dictionary:      dict,
		Articles:        articles,
		LetterFilter:    letterFilter,
		Pagination: htmlui.Pagination{
			Current: urlQuery.Get("page").(*htmlui.IntegerQueryParam).Value(),
			Total:   int(math.Ceil(float64(total) / pageSize)),
			PageToURL: func(n int) string {
				urlQuery := urlQuery.Clone()
				urlQuery.Get("page").(*htmlui.IntegerQueryParam).SetValue(n)
				result := urlQuery.Encode()
				if result != "" {
					result = "?" + result
				}
				return rctx.R.URL.Path + result
			},
		},
	})
	if err != nil {
		return fmt.Errorf("render html: %w", err)
	}

	return nil
}

func calcPages(c, total int) []int {
	result := []int{}

	return result
}
