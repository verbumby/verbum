package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/verbumby/verbum/backend/chttp"
	"github.com/verbumby/verbum/backend/dictionary"
)

func APIDictionaryAbbrevs(w http.ResponseWriter, rctx *chttp.Context) error {
	d := dictionary.Get(chi.URLParam(rctx.R, "dictionary"))
	if d == nil || d.Abbrevs() == nil {
		return APINotFound(w, rctx)
	}

	if err := json.NewEncoder(w).Encode(d.Abbrevs()); err != nil {
		return fmt.Errorf("encode response: %w", err)
	}
	return nil
}
