package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"gopkg.in/reform.v1"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/verbumby/verbum/backend/pkg/app"
	"github.com/verbumby/verbum/backend/pkg/article"
	"github.com/verbumby/verbum/backend/pkg/chttp"
	"github.com/verbumby/verbum/backend/pkg/db"
	"github.com/verbumby/verbum/backend/pkg/dict"
	"github.com/verbumby/verbum/backend/pkg/fts"
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
			name:    "admin",
			files:   []string{"./templates/admin/layout.html"},
			funcMap: template.FuncMap{},
		},
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
	r.Handle("/admin/api/dictionaries", chttp.MakeHandler(
		(&RecordListHandler{
			Table: dict.DictTable,
			DB:    db.DB,
		}).ServeHTTP,
		chttp.AuthMiddleware,
	)).Methods(http.MethodGet)
	r.Handle("/admin/api/dictionaries", chttp.MakeHandler(
		(&RecordSaveHandler{
			Table: dict.DictTable,
			DB:    db.DB,
		}).ServeHTTP,
		chttp.AuthMiddleware,
	)).Methods(http.MethodPost)
	r.Handle("/admin/api/articles", chttp.MakeHandler(
		(&RecordListHandler{
			Table: article.ArticleTable,
			DB:    db.DB,
			Filters: []app.Filter{
				&article.FilterDictID{},
				&article.FilterTitlePrefix{},
			},
		}).ServeHTTP,
		chttp.AuthMiddleware,
	)).Methods(http.MethodGet)
	r.Handle("/admin/api/articles", chttp.MakeHandler(
		(&article.RecordSaveHandler{
			ModelMeta: article.ArticleMeta,
			DB:        db.DB,
		}).ServeHTTP,
		chttp.AuthMiddleware,
	)).Methods(http.MethodPost)
	r.HandleFunc("/admin/api/articles/{ID}", chttp.MakeHandler(
		(&RecordFetchHandler{
			ModelMeta: article.ArticleMeta,
			DB:        db.DB,
		}).ServeHTTP,
		chttp.AuthMiddleware,
	)).Methods(http.MethodGet)
	r.HandleFunc("/admin/auth", chttp.MakeHandler(chttp.AuthHandler))
	r.PathPrefix("/admin/").HandlerFunc(chttp.MakeHandler(IndexHandler, chttp.AuthMiddleware))

	statics := http.FileServer(http.Dir("statics"))
	r.PathPrefix("/statics/").Handler(http.StripPrefix("/statics/", statics))
	r.HandleFunc("/_suggest", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		query := "SELECT typeahead, MAX(WEIGHT()) mwt " +
			"FROM typeaheads " +
			"WHERE MATCH(?) " +
			"GROUP BY typeahead " +
			"ORDER BY mwt DESC " +
			"LIMIT 10 " +
			"OPTION ranker=sph04 "

		rows, err := fts.Sphinx.Query(query, q)
		if err != nil {
			log.Println(errors.Wrap(err, "find typeaheads"))
		}
		data := []string{}
		for rows.Next() {
			var th string
			var mwt int32
			if err := rows.Scan(&th, &mwt); err != nil {
				log.Println(errors.Wrap(err, "scan typeaheads"))
			}

			data = append(data, th)
		}

		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Println(errors.Wrap(err, "encode response"))
		}
	})
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")

		var articles []article.Article
		var err error
		if q != "" {
			articles, err = func() ([]article.Article, error) {
				rows, err := fts.Sphinx.Query(
					"SELECT article_id, MAX(WEIGHT()) mw "+
						"FROM headwords "+
						"WHERE MATCH(?) "+
						"GROUP BY article_id "+
						"ORDER BY mw DESC "+
						"LIMIT 20 "+
						"OPTION ranker=sph04",
					q,
				)
				if err != nil {
					return nil, errors.Wrap(err, "query sphinx")
				}
				defer rows.Close()

				articleIDs := []int32{}
				for rows.Next() {
					var articleID int32
					var maxWeight float32
					if err := rows.Scan(&articleID, &maxWeight); err != nil {
						return nil, errors.Wrap(err, "sphinx rows scan")
					}
					articleIDs = append(articleIDs, articleID)
				}

				result := make([]article.Article, len(articleIDs))
				for i, id := range articleIDs {
					if err := db.DB.FindByPrimaryKeyTo(&result[i], id); err != nil {
						return nil, errors.Wrapf(err, "find article by pk %d", id)
					}
				}
				return result, nil
			}()
			if err != nil {
				log.Println(err)
			}
		}
		tm.Render("index", w, struct {
			Articles []article.Article
			Q        string
		}{
			Articles: articles,
			Q:        q,
		})
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
