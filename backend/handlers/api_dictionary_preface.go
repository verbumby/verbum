package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/verbumby/verbum/backend/chttp"
	"github.com/verbumby/verbum/backend/dictionary"
)

func APIDictionaryPreface(w http.ResponseWriter, rctx *chttp.Context) error {
	d := dictionary.Get(chi.URLParam(rctx.R, "dictionary"))
	if d == nil || d.Preface() == "" {
		return APINotFound(w, rctx)
	}

	if err := json.NewEncoder(w).Encode(d.Preface()); err != nil {
		return fmt.Errorf("encode response: %w", err)
	}
	return nil
}
