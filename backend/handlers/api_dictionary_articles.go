package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/verbumby/verbum/backend/article"
	"github.com/verbumby/verbum/backend/chttp"
	"github.com/verbumby/verbum/backend/dictionary"
	"github.com/verbumby/verbum/backend/htmlui"
)

// APIDictionaryArticles handles dictionary articles request
func APIDictionaryArticles(w http.ResponseWriter, rctx *chttp.Context) error {
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

	const pageSize = 10
	prefixMusts := []interface{}{}
	for i, p := range prefix {
		prefixMusts = append(prefixMusts, map[string]interface{}{
			"term": map[string]interface{}{
				fmt.Sprintf("Prefix.Letter%d", i+1): string(p),
			},
		})
	}
	articles, total, err := article.Query("/dict-"+d.IndexID()+"/_search", map[string]interface{}{
		"track_total_hits": true,
		"from":             (urlQuery.Get("page").(*htmlui.IntegerQueryParam).Value() - 1) * pageSize,
		"size":             pageSize,
		"sort": []any{
			"SortKey",
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

	type articleview struct {
		ID           string
		Content      string
		DictionaryID string
	}

	articleviews := []articleview{}
	for _, a := range articles {
		articleviews = append(articleviews, articleview{
			ID:           a.ID,
			Content:      string(a.Dictionary.ToHTML(a.Content)),
			DictionaryID: a.Dictionary.ID(),
		})
	}

	type paginationview struct {
		Current int
		Total   int
	}

	if err := json.NewEncoder(w).Encode(struct {
		DictID     string
		Prefix     string
		Articles   []articleview
		Pagination paginationview
	}{
		DictID:   d.ID(),
		Prefix:   string(prefix),
		Articles: articleviews,
		Pagination: paginationview{
			Current: urlQuery.Get("page").(*htmlui.IntegerQueryParam).Value(),
			Total:   int(math.Ceil(float64(total) / pageSize)),
		},
	}); err != nil {
		return fmt.Errorf("encode response: %w", err)
	}

	return nil
}
