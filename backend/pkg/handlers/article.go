package handlers

import (
	"net/http"

	"github.com/verbumby/verbum/backend/pkg/htmlui"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/verbumby/verbum/backend/pkg/article"
	"github.com/verbumby/verbum/backend/pkg/chttp"
	"github.com/verbumby/verbum/backend/pkg/tm"
)

// Article handles article page request
func Article(w http.ResponseWriter, rctx *chttp.Context) error {
	vars := mux.Vars(rctx.R)
	dID := vars["dictionary"]
	aID := vars["article"]

	a, err := article.Get(dID, aID)
	if err != nil {
		return errors.Wrap(err, "get article")
	}

	err = tm.Render("article", w, struct {
		PageTitle       string
		PageDescription string
		MetaRobotsTag   htmlui.MetaRobotsTag
		Article         article.Article
	}{
		PageTitle:       a.Title + " - " + a.Dictionary.Title(),
		PageDescription: a.Title + " - " + a.Dictionary.Title(),
		MetaRobotsTag:   htmlui.MetaRobotsTag{Index: true, Follow: false},
		Article:         a,
	})
	if err != nil {
		return errors.Wrap(err, "render html")
	}

	return nil
}
