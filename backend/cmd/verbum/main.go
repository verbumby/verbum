package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/verbumby/verbum/backend/pkg/article"
	"github.com/verbumby/verbum/backend/pkg/dict"
	"github.com/verbumby/verbum/backend/pkg/fts"
	"github.com/verbumby/verbum/backend/pkg/tm"
)

var (
	// DB reform database handler
	DB *reform.DB

	// Config global application config
	Config struct {
		DBHost     string
		SphinxHost string
	}
)

func main() {
	err := bootstrap()
	if err != nil {
		log.Fatal(err)
	}
}

func bootstrap() error {
	flag.StringVar(&Config.DBHost, "db-host", "localhost", "hostname of the database server")
	flag.StringVar(&Config.SphinxHost, "sphinx-host", "localhost", "hostname of the sphinx server")
	flag.Parse()

	// TODO: parametrize db connection
	db, err := sql.Open("mysql", fmt.Sprintf("root@tcp(%s:3306)/verbum", Config.DBHost))
	if err != nil {
		return errors.Wrap(err, "open db")
	}
	DB = reform.NewDB(db, mysql.Dialect, reform.NewPrintfLogger(log.Printf))

	sphinxConnString := fmt.Sprintf("tcp(%s:9306)/?interpolateParams=true", Config.SphinxHost)
	if err := fts.Initialize(sphinxConnString); err != nil {
		return errors.Wrap(err, "fts initialize")
	}

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
					return DB.FindByPrimaryKeyFrom(dict.DictTable, id)
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
	r.Handle("/admin/api/dictionaries", &RecordListHandler{
		Table: dict.DictTable,
		DB:    DB,
	}).Methods(http.MethodGet)
	r.Handle("/admin/api/dictionaries", &RecordSaveHandler{
		Table: dict.DictTable,
		DB:    DB,
	}).Methods(http.MethodPost)
	r.Handle("/admin/api/articles", &RecordListHandler{
		Table: article.ArticleTable,
		DB:    DB,
	}).Methods(http.MethodGet)
	r.Handle("/admin/api/articles", &RecordSaveHandler{
		Table:     article.ArticleTable,
		DB:        DB,
		AfterSave: article.Index,
	}).Methods(http.MethodPost)
	r.Handle("/admin/api/articles/{ID}", &RecordFetchHandler{
		Table: article.ArticleTable,
		DB:    DB,
	}).Methods(http.MethodGet)
	r.PathPrefix("/admin/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dicts, err := DB.SelectAllFrom(dict.DictTable, "")
		if err != nil {
			log.Println(errors.Wrap(err, "select all dicts"))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		data := struct {
			Dicts []reform.Struct
		}{
			Dicts: dicts,
		}
		if err := tm.Render("admin", w, data); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}
	})

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
					if err := DB.FindByPrimaryKeyTo(&result[i], id); err != nil {
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

	log.Println("listening on :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		return errors.Wrap(err, "listen and serve")
	}

	return nil
}
