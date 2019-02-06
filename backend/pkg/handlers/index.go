package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/verbumby/verbum/backend/pkg/chttp"
	"github.com/verbumby/verbum/backend/pkg/tm"
)

// Index handles / request
func Index(w http.ResponseWriter, rctx *chttp.Context) error {
	pageTitle := "Verbum - Анлайн Слоўнік Беларускай Мовы"
	pageDescription := pageTitle
	q := rctx.R.URL.Query().Get("q")

	type articleView struct {
		DictTitle string
		Content   string
	}
	var articles []articleView
	var err error
	if q != "" {
		pageTitle = q + " - Пошук"
		qbytes, err := json.Marshal(q)
		if err != nil {
			return errors.Wrap(err, "marshal q")
		}

		query := `{
			"query": {
				"multi_match": {
					"query": ` + string(qbytes) + `,
					"fields": [
						"Headword^4",
						"Headword.Smaller^3",
						"HeadwordAlt^2",
						"HeadwordAlt.Smaller^1"
					]
				}
			}
		}`

		url := viper.GetString("elastic.addr") + "/dict-*/_search"
		resp, err := http.Post(url, "application/json", strings.NewReader(query))
		if err != nil {
			return errors.Wrap(err, "query elastic")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respbytes, _ := ioutil.ReadAll(resp.Body)
			return fmt.Errorf(
				"query elastic: expected %d, got %d: %s",
				http.StatusOK,
				resp.StatusCode,
				string(respbytes),
			)
		}

		respdata := struct {
			Hits struct {
				Hits []struct {
					Source struct {
						Content string
					} `json:"_source"`
					Index string `json:"_index"`
				} `json:"hits"`
			} `json:"hits"`
		}{}
		if err := json.NewDecoder(resp.Body).Decode(&respdata); err != nil {
			return errors.Wrap(err, "decode elastic resp")
		}

		dicts := map[string]string{}
		for _, hit := range respdata.Hits.Hits {
			dictID := strings.TrimPrefix(hit.Index, "dict-")
			if _, ok := dicts[dictID]; !ok {
				url := viper.GetString("elastic.addr") + "/dicts/_doc/" + dictID
				resp, err := http.Get(url)
				if err != nil {
					return errors.Wrapf(err, "query dict %s: new request", dictID)
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					respbytes, _ := ioutil.ReadAll(resp.Body)
					return fmt.Errorf(
						"query dict %s: expected %d, got %d: %s",
						dictID,
						http.StatusOK,
						resp.StatusCode,
						string(respbytes),
					)
				}

				respdata := struct {
					Source struct {
						Title string
					} `json:"_source"`
				}{}
				if err := json.NewDecoder(resp.Body).Decode(&respdata); err != nil {
					return errors.Wrapf(err, "query dict %s: decode elastic resp", dictID)
				}

				dicts[dictID] = respdata.Source.Title
			}

			articles = append(articles, articleView{
				DictTitle: dicts[dictID],
				Content:   hit.Source.Content,
			})
		}
	}

	if len(articles) > 0 {
		pageDescription = articles[0].Content
	}
	err = tm.Render("index", w, struct {
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
