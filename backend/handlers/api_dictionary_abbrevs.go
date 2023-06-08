package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/verbumby/verbum/backend/chttp"
	"github.com/verbumby/verbum/backend/dictionary"
)

func APIDictionaryAbbrevs(w http.ResponseWriter, rctx *chttp.Context) error {
	vars := mux.Vars(rctx.R)

	d := dictionary.Get(vars["dictionary"])
	if d == nil || d.Abbrevs() == nil {
		return APINotFound(w, rctx)
	}

	if err := json.NewEncoder(w).Encode(d.Abbrevs()); err != nil {
		return fmt.Errorf("encode response: %w", err)
	}
	return nil
}
