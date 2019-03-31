package handlers

import (
	"math"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/verbumby/verbum/backend/pkg/article"
	"github.com/verbumby/verbum/backend/pkg/chttp"
	"github.com/verbumby/verbum/backend/pkg/dictionary"
	"github.com/verbumby/verbum/backend/pkg/htmlui"
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

	urlQuery := htmlui.Query([]htmlui.QueryParam{
		htmlui.NewIntegerQueryParam("page", 1),
	})
	urlQuery.From(rctx.R.URL.Query())

	articles, total, err := article.Query("/dict-"+dictID+"/_search", map[string]interface{}{
		"from": (urlQuery.Get("page").(*htmlui.IntegerQueryParam).Value() - 1) * 10,
		"size": 10,
		"sort": []interface{}{
			"Title",
		},
	})
	if err != nil {
		return errors.Wrap(err, "query articles")
	}

	err = tm.Render("dictionary", w, struct {
		PageTitle       string
		PageDescription string
		Dictionary      dictionary.Dictionary
		Articles        []article.Article
		Pagination      htmlui.Pagination
	}{
		PageTitle:       dict.Title,
		PageDescription: dict.Title,
		Dictionary:      dict,
		Articles:        articles,
		Pagination: htmlui.Pagination{
			Current: urlQuery.Get("page").(*htmlui.IntegerQueryParam).Value(),
			Total:   int(math.Ceil(float64(total) / 10)),
			PageToURL: func(n int) string {
				urlQuery := urlQuery.Clone()
				urlQuery.Get("page").(*htmlui.IntegerQueryParam).SetValue(n)
				return "?" + urlQuery.Encode()
			},
		},
	})
	if err != nil {
		return errors.Wrap(err, "render html")
	}

	return nil
}

func calcPages(c, total int) []int {
	result := []int{}

	return result
}
