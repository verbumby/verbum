package handlers

import (
	"math"
	"net/http"
	"strconv"

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

	urlQuery := htmlui.URLQuery([]htmlui.URLQueryParam{
		htmlui.NewIntegerURLQueryParam("page", 1),
	}).From(rctx.R.URL.Query())

	articles, total, err := article.Query("/dict-"+dictID+"/_search", map[string]interface{}{
		"from": (urlQuery.Get("page").(*htmlui.IntegerURLQueryParam).Value() - 1) * 10,
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
			Current: urlQuery.Get("page").(*htmlui.IntegerURLQueryParam).Value(),
			Total:   int(math.Ceil(float64(total) / 10)),
			PageToURL: func(n int) string {
				return "?" + urlQuery.With("page", strconv.FormatInt(int64(n), 10)).Encode()
			},
		},
		// PaginationPages: calcPages(12,),
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
