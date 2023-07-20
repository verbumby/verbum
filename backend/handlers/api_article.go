package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/verbumby/verbum/backend/article"
	"github.com/verbumby/verbum/backend/chttp"
	"github.com/verbumby/verbum/backend/dictionary"
)

// APIArticle handles article page request
func APIArticle(w http.ResponseWriter, rctx *chttp.Context) error {
	d := dictionary.Get(chi.URLParam(rctx.R, "dictionary"))
	if d == nil {
		return APINotFound(w, rctx)
	}

	aID := chi.URLParam(rctx.R, "article")
	var err error
	aID, err = url.QueryUnescape(aID)
	if err != nil {
		return fmt.Errorf("unescape aID: %w", err)
	}

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
