package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/verbumby/verbum/dict"
	"github.com/verbumby/verbum/tm"
)

var (
	DB *reform.DB
)

func main() {
	err := bootstrap()
	if err != nil {
		log.Fatal(err)
	}
}

func bootstrap() error {

	// TODO: parametrize db connection
	db, err := sql.Open("mysql", "root@/verbum")
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

	http.HandleFunc("/admin/api/dictionaries", dictionaryCollectionHandler)
	http.HandleFunc("/admin/", func(w http.ResponseWriter, r *http.Request) {
		if err := tm.Render("admin", w, nil); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}
	})
	statics := http.FileServer(http.Dir("statics"))
	http.Handle("/statics/", http.StripPrefix("/statics/", statics))

	log.Println("listening on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return errors.Wrap(err, "listen and serve")
	}

	return nil
}

func dictionaryCollectionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		d := &dict.Dict{}
		if err := json.NewDecoder(r.Body).Decode(d); err != nil {
			http.Error(w, errors.Wrap(err, "decode body").Error(), http.StatusBadRequest)
			return
		}

		if err := DB.Save(d); err != nil {
			log.Println(errors.Wrap(err, "save dictionary"))
			http.Error(w, "", http.StatusInternalServerError)
		}
		log.Println(d)
	}
}
