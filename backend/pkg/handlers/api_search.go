package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/verbumby/verbum/backend/pkg/article"
	"github.com/verbumby/verbum/backend/pkg/chttp"
)

// APISearch search endpoint
func APISearch(w http.ResponseWriter, rctx *chttp.Context) error {
	var err error
	q := rctx.R.URL.Query().Get("q")
	if len(q) > 1000 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil
	}

	reqbody := map[string]interface{}{
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
	}

	articles, _, err := article.Query("/dict-*/_search", reqbody)
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

	if err := json.NewEncoder(w).Encode(struct {
		Articles []articleview
	}{
		Articles: articleviews,
	}); err != nil {
		return fmt.Errorf("encode response: %w", err)
	}

	return nil
}
