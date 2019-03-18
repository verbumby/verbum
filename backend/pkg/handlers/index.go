package handlers

import (
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/verbumby/verbum/backend/pkg/chttp"
	"github.com/verbumby/verbum/backend/pkg/storage"
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

	err := tm.Render("index", w, struct {
		Q               string
		PageTitle       string
		PageDescription string
	}{
		Q:               "",
		PageTitle:       pageTitle,
		PageDescription: pageDescription,
	})
	if err != nil {
		return errors.Wrap(err, "render html")
	}
	return nil
}

func search(w http.ResponseWriter, rctx *chttp.Context) error {
	type articleView struct {
		DictTitle string
		Content   string
	}
	var articles []articleView
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

	respbody := struct {
		Hits struct {
			Hits []struct {
				Source struct {
					Content string
				} `json:"_source"`
				Index string `json:"_index"`
			} `json:"hits"`
		} `json:"hits"`
	}{}
	if err := storage.Post("/dict-*/_search", reqbody, &respbody); err != nil {
		return errors.Wrap(err, "query elastic")
	}

	dicts := map[string]string{}
	for _, hit := range respbody.Hits.Hits {
		dictID := strings.TrimPrefix(hit.Index, "dict-")
		if _, ok := dicts[dictID]; !ok {
			respbody := struct {
				Source struct {
					Title string
				} `json:"_source"`
			}{}

			if err := storage.Get("/dicts/_doc/"+dictID, &respbody); err != nil {
				return errors.Wrapf(err, "query dict %s", dictID)
			}

			dicts[dictID] = respbody.Source.Title
		}

		articles = append(articles, articleView{
			DictTitle: dicts[dictID],
			Content:   hit.Source.Content,
		})
	}

	if len(articles) > 0 {
		pageDescription = articles[0].Content
	}
	err = tm.Render("search-results", w, struct {
		Articles        []articleView
		Q               string
		PageTitle       string
		PageDescription string
	}{
		Articles:        articles,
		Q:               q,
		PageTitle:       pageTitle,
		PageDescription: pageDescription,
	})
	if err != nil {
		return errors.Wrap(err, "render html")
	}
	return nil
}
