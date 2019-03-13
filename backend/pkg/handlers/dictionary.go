package handlers

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/verbumby/verbum/backend/pkg/chttp"
	"github.com/verbumby/verbum/backend/pkg/tm"
)

// Dictionary handles / request
func Dictionary(w http.ResponseWriter, rctx *chttp.Context) error {
	err := tm.Render("dictionary", w, struct {
		PageTitle       string
		PageDescription string
	}{
		PageTitle:       "Page TItle goes here",
		PageDescription: "",
	})
	if err != nil {
		return errors.Wrap(err, "render html")
	}

	return nil
}
