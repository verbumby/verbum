package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/verbumby/verbum/backend/pkg/article"
	"github.com/verbumby/verbum/backend/pkg/chttp"
	"github.com/verbumby/verbum/backend/pkg/dictionary"
	"github.com/verbumby/verbum/backend/pkg/tm"
)

// Dictionary handles dictionary page request
func Dictionary(w http.ResponseWriter, rctx *chttp.Context) error {
	vars := mux.Vars(rctx.R)
	dictID := vars["dictionary"]

	dict, err := dictionary.Get(dictID)
	if err != nil {
		return errors.Wrapf(err, "dictionary get %s", dictID)
	}

	articles, err := article.Query("/dict-"+dictID+"/_search", map[string]interface{}{
		"size": 20,
	})
	if err != nil {
		return errors.Wrap(err, "query articles")
	}

	err = tm.Render("dictionary", w, struct {
		PageTitle       string
		PageDescription string
		Dictionary      dictionary.Dictionary
		Articles        []article.Article
	}{
		PageTitle:       dict.Title,
		PageDescription: dict.Title,
		Dictionary:      dict,
		Articles:        articles,
	})
	if err != nil {
		return errors.Wrap(err, "render html")
	}

	return nil
}
