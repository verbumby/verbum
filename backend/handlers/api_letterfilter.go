package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/verbumby/verbum/backend/chttp"
	"github.com/verbumby/verbum/backend/dictionary"
	"github.com/verbumby/verbum/backend/htmlui"
	"github.com/verbumby/verbum/backend/storage"
)

// APILetterFilter handle letter filter request
func APILetterFilter(w http.ResponseWriter, rctx *chttp.Context) error {
	d := dictionary.Get(chi.URLParam(rctx.R, "dictionary"))
	if d == nil {
		return APINotFound(w, rctx)
	}

	urlQuery := htmlui.Query([]htmlui.QueryParam{
		htmlui.NewIntegerQueryParam("page", 1),
		htmlui.NewStringQueryParam("prefix", ""),
	})
	urlQuery.From(rctx.R.URL.Query())

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
	if err := storage.Post("/dict-"+d.IndexID()+"/_search", aggsreqbody, &aggsrespbody); err != nil {
		return fmt.Errorf("aggs query: %w", err)
	}

	letterFilter := htmlui.LetterFilter{
		Prefix: prefix,
		LetterLink: func(prefix string) string {
			return prefix
		},
	}
	letterFilter.AddLevel(aggsrespbody.Aggregations.Prefix.Letter1.Buckets)
	if aggsrespbody.Aggregations.Prefix.Letter2 != nil {
		letterFilter.AddLevel(aggsrespbody.Aggregations.Prefix.Letter2.Letter2.Buckets)
	}
	if aggsrespbody.Aggregations.Prefix.Letter3 != nil {
		letterFilter.AddLevel(aggsrespbody.Aggregations.Prefix.Letter3.Letter3.Letter3.Buckets)
	}

	type letterfilterview struct {
		DictID  string
		Prefix  string
		Entries [][]htmlui.LetterFilterLink
	}

	if err := json.NewEncoder(w).Encode(letterfilterview{
		DictID:  d.ID(),
		Prefix:  string(prefix),
		Entries: letterFilter.Links(),
	}); err != nil {
		return fmt.Errorf("encode response: %w", err)
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
