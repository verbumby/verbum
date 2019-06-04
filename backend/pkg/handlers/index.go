package handlers

import (
	"net/http"

	"github.com/verbumby/verbum/backend/pkg/dictionary"

	"github.com/pkg/errors"
	"github.com/verbumby/verbum/backend/pkg/article"
	"github.com/verbumby/verbum/backend/pkg/chttp"
	"github.com/verbumby/verbum/backend/pkg/htmlui"
	"github.com/verbumby/verbum/backend/pkg/tm"
)

// Index handles / request
func Index(w http.ResponseWriter, rctx *chttp.Context) error {
	q := rctx.R.URL.Query().Get("q")
	if q != "" {
		return search(w, rctx)
	}
	return index(w, rctx)

}

func index(w http.ResponseWriter, rctx *chttp.Context) error {
	pageTitle := "Verbum - Анлайн Слоўнік Беларускай Мовы"
	pageDescription := pageTitle

	dicts := dictionary.GetAll()

	err := tm.Render("index", w, struct {
		Q               string
		PageTitle       string
		PageDescription string
		MetaRobotsTag   htmlui.MetaRobotsTag
		Dictionaries    []dictionary.Dictionary
	}{
		Q:               "",
		PageTitle:       pageTitle,
		PageDescription: pageDescription,
		MetaRobotsTag:   htmlui.MetaRobotsTag{Index: true, Follow: true},
		Dictionaries:    dicts,
	})
	if err != nil {
		return errors.Wrap(err, "render html")
	}
	return nil
}

func search(w http.ResponseWriter, rctx *chttp.Context) error {
	var err error
	q := rctx.R.URL.Query().Get("q")
	pageTitle := q + " - Пошук"
	pageDescription := pageTitle

	reqbody := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": q,
				"fields": []string{
					"Headword^4",
					"Headword.Smaller^3",
					"HeadwordAlt^2",
					"HeadwordAlt.Smaller^1",
				},
			},
		},
	}

	articles, _, err := article.Query("/dict-*/_search", reqbody)
	if err != nil {
		return errors.Wrap(err, "query articles")
	}

	err = tm.Render("search-results", w, struct {
		Articles        []article.Article
		Q               string
		PageTitle       string
		PageDescription string
		MetaRobotsTag   htmlui.MetaRobotsTag
	}{
		Articles:        articles,
		Q:               q,
		PageTitle:       pageTitle,
		PageDescription: pageDescription,
		MetaRobotsTag:   htmlui.MetaRobotsTag{Index: false, Follow: false},
	})
	if err != nil {
		return errors.Wrap(err, "render html")
	}
	return nil
}
