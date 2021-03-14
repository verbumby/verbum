package handlers

import (
	"net/http"

	"github.com/verbumby/verbum/backend/pkg/chttp"
)

// APINotFound returns 404 to the caller
func APINotFound(w http.ResponseWriter, rctx *chttp.Context) error {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	return nil
}
