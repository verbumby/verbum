package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/verbumby/verbum/backend/pkg/article"
	"github.com/verbumby/verbum/backend/pkg/chttp"
	"github.com/verbumby/verbum/backend/pkg/dictionary"
)

// APIArticle handles article page request
func APIArticle(w http.ResponseWriter, rctx *chttp.Context) error {
	vars := mux.Vars(rctx.R)
	d := dictionary.Get(vars["dictionary"])
	if d == nil {
		return APINotFound(w, rctx)
	}

	aID := vars["article"]

	a, err := article.Get(d, aID)
	if err != nil {
		return fmt.Errorf("get article: %w", err)
	}

	if a.ID == "" {
		return APINotFound(w, rctx)
	}

	type articleview struct {
		ID           string
		Title        string
		Headword     []string
		Content      string
		DictionaryID string
	}
	resp := articleview{
		ID:           a.ID,
		Title:        a.Title,
		Headword:     a.Headword,
		Content:      string(a.Dictionary.ToHTML(a.Content)),
		DictionaryID: a.Dictionary.ID(),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return fmt.Errorf("encode response: %w", err)
	}

	return nil
}
