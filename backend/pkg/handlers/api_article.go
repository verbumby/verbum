package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/verbumby/verbum/backend/pkg/article"
	"github.com/verbumby/verbum/backend/pkg/chttp"
)

// APIArticle handles article page request
func APIArticle(w http.ResponseWriter, rctx *chttp.Context) error {
	vars := mux.Vars(rctx.R)
	dID := vars["dictionary"]
	aID := vars["article"]

	a, err := article.Get(dID, aID)
	if err != nil {
		return fmt.Errorf("get article: %w", err)
	}

	type articleview struct {
		ID           string
		Title        string
		Content      string
		DictionaryID string
	}
	resp := articleview{
		ID:           a.ID,
		Title:        a.Title,
		Content:      string(a.Dictionary.ToHTML(a.Content, a.Title)),
		DictionaryID: a.Dictionary.ID(),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return fmt.Errorf("encode response: %w", err)
	}

	return nil
}
