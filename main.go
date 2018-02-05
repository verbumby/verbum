package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/verbumby/verbum/article"
	"github.com/verbumby/verbum/dict"
	"github.com/verbumby/verbum/tm"
)

var (
	// DB reform database handler
	DB *reform.DB

	// Config global application config
	Config struct {
		DBHost string
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
	flag.Parse()

	// TODO: parametrize db connection
	db, err := sql.Open("mysql", fmt.Sprintf("root@tcp(%s:3306)/verbum", Config.DBHost))
	if err != nil {
		return errors.Wrap(err, "open db")
	}
	DB = reform.NewDB(db, mysql.Dialect, reform.NewPrintfLogger(log.Printf))

	templates := map[string][]string{
		"admin": []string{"./templates/admin/layout.html"},
	}
	for tn, files := range templates {
		if err := tm.Compile(tn, files); err != nil {
			return errors.Wrapf(err, "compile %s template", tn)
		}
	}

	r := mux.NewRouter()
	r.Handle("/admin/api/dictionaries", &RecordListHandler{
		Table: dict.DictTable,
		DB:    DB,
	}).Methods(http.MethodGet)
	r.Handle("/admin/api/dictionaries", &RecordCreateHandler{
		Table: dict.DictTable,
		DB:    DB,
	}).Methods(http.MethodPost)
	r.Handle("/admin/api/articles", &RecordListHandler{
		Table: article.ArticleTable,
		DB:    DB,
	}).Methods(http.MethodGet)
	r.Handle("/admin/api/articles", &RecordCreateHandler{
		Table: article.ArticleTable,
		DB:    DB,
	}).Methods(http.MethodPost)
	r.PathPrefix("/admin/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := tm.Render("admin", w, nil); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}
	})

	statics := http.FileServer(http.Dir("statics"))
	r.PathPrefix("/statics/").Handler(http.StripPrefix("/statics/", statics))

	log.Println("listening on :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		return errors.Wrap(err, "listen and serve")
	}

	return nil
}
