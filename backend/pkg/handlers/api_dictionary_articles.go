package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/verbumby/verbum/backend/pkg/article"
	"github.com/verbumby/verbum/backend/pkg/chttp"
	"github.com/verbumby/verbum/backend/pkg/dictionary"
	"github.com/verbumby/verbum/backend/pkg/htmlui"
)

// APIDictionaryArticles handles dictionary articles request
func APIDictionaryArticles(w http.ResponseWriter, rctx *chttp.Context) error {
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
	articles, total, err := article.Query("/dict-"+dictID+"/_search", map[string]interface{}{
		"track_total_hits": true,
		"from":             (urlQuery.Get("page").(*htmlui.IntegerQueryParam).Value() - 1) * pageSize,
		"size":             pageSize,
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
		DictID:   dictID,
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
