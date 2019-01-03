package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	reform "gopkg.in/reform.v1"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/verbumby/verbum/backend/pkg/app"
	"github.com/verbumby/verbum/backend/pkg/chttp"
	"github.com/verbumby/verbum/backend/pkg/db"
	"github.com/verbumby/verbum/backend/pkg/dict"
	"github.com/verbumby/verbum/backend/pkg/tm"
)

func main() {
	if err := app.Bootstrap(); err != nil {
		log.Fatal(err)
	}

	if err := bootstrapServer(); err != nil {
		log.Fatal(err)
	}
}

func bootstrapServer() error {
	templates := []struct {
		name    string
		files   []string
		funcMap template.FuncMap
	}{
		{
			name:  "index",
			files: []string{"./templates/index.html"},
			funcMap: template.FuncMap{
				"dictByPK": func(id int32) (reform.Record, error) {
					return db.DB.FindByPrimaryKeyFrom(dict.DictTable, id)
				},
			},
		},
	}

	for _, t := range templates {
		if err := tm.Compile(t.name, t.files, t.funcMap); err != nil {
			return errors.Wrapf(err, "compile %s template", t.name)
		}
	}

	r := mux.NewRouter()
	statics := http.FileServer(http.Dir("statics"))
	r.PathPrefix("/statics/").Handler(http.StripPrefix("/statics/", statics))
	r.HandleFunc("/_suggest", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")

		qbytes, err := json.Marshal(q)
		if err != nil {
			log.Println(errors.Wrap(err, "marshal q"))
			return
		}
		query := `{
			"_source": false,
			"suggest": {
				"HeadwordSuggest": {
					"prefix": ` + string(qbytes) + `,
					"completion": {
						"field": "Suggest",
						"skip_duplicates": true,
						"size": 10
					}
				}
			}
		}`

		url := viper.GetString("elastic.addr") + "/dict-*/_search"
		resp, err := http.Post(url, "application/json", strings.NewReader(query))
		if err != nil {
			log.Println(errors.Wrap(err, "query elastic"))
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respbytes, _ := ioutil.ReadAll(resp.Body)
			log.Println(fmt.Errorf("query elastic: expected %d, got %d: %s", http.StatusOK, resp.StatusCode, string(respbytes)))
			return
		}

		respdata := struct {
			Suggest struct {
				HeadwordSuggest []struct {
					Options []struct {
						Text string `json:"text"`
					} `json:"options"`
				}
			} `json:"suggest"`
		}{}
		if err := json.NewDecoder(resp.Body).Decode(&respdata); err != nil {
			log.Println(errors.Wrap(err, "unmarshal elastic resp"))
			return
		}

		data := []string{}
		for _, hws := range respdata.Suggest.HeadwordSuggest {
			for _, opt := range hws.Options {
				data = append(data, opt.Text)
			}
		}

		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Println(errors.Wrap(err, "encode response"))
		}
	})
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pageTitle := "Verbum - Анлайн Слоўнік Беларускай Мовы"
		pageDescription := pageTitle
		q := r.URL.Query().Get("q")

		var articles []string
		var err error
		if q != "" {
			pageTitle = q + " - Пошук"
			qbytes, err := json.Marshal(q)
			if err != nil {
				log.Println(errors.Wrap(err, "marshal q"))
				return
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
				log.Println(errors.Wrap(err, "query elastic"))
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				respbytes, _ := ioutil.ReadAll(resp.Body)
				log.Println(fmt.Errorf("query elastic: expected %d, got %d: %s", http.StatusOK, resp.StatusCode, string(respbytes)))
				return
			}

			respdata := struct {
				Hits struct {
					Hits []struct {
						Source struct {
							Content string
						} `json:"_source"`
					} `json:"hits"`
				} `json:"hits"`
			}{}
			if err := json.NewDecoder(resp.Body).Decode(&respdata); err != nil {
				log.Println(errors.Wrap(err, "decode elastic resp"))
				return
			}

			for _, hit := range respdata.Hits.Hits {
				articles = append(articles, hit.Source.Content)
			}
		}

		if len(articles) > 0 {
			pageDescription = articles[0]
		}
		err = tm.Render("index", w, struct {
			Articles        []string
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
			log.Println(err)
		}
	})

	chttp.InitCookieManager()

	if viper.IsSet("http.addr") {
		go func() {
			statics := http.FileServer(http.Dir(viper.GetString("http.acmeChallengeRoot")))
			r := http.NewServeMux()
			r.Handle("/.well-known/acme-challenge/", http.StripPrefix("/.well-known/acme-challenge/", statics))
			r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
				if req.Method == http.MethodGet {
					target := "https://" + req.Host + req.URL.Path
					if len(req.URL.RawQuery) > 0 {
						target += "?" + req.URL.RawQuery
					}
					http.Redirect(w, req, target, http.StatusTemporaryRedirect)
				} else {
					http.NotFound(w, req)
				}
			})
			http.ListenAndServe(viper.GetString("http.addr"), r)
		}()
	}

	log.Printf("listening on %s", viper.GetString("https.addr"))
	err := http.ListenAndServeTLS(
		viper.GetString("https.addr"),
		viper.GetString("https.certFile"),
		viper.GetString("https.keyFile"),
		r,
	)
	if err != nil {
		return errors.Wrap(err, "listen and serve")
	}

	return nil
}
