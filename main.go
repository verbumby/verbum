package main

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/verbumby/verbum/tm"
)

func main() {
	err := bootstrap()
	if err != nil {
		log.Fatal(err)
	}
}

func bootstrap() error {
	templates := map[string][]string{
		"admin": []string{"./templates/admin/layout.html"},
	}
	for tn, files := range templates {
		if err := tm.Compile(tn, files); err != nil {
			return errors.Wrapf(err, "compile %s template", tn)
		}
	}

	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		if err := tm.Render("admin", w, nil); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}
	})

	log.Println("listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return errors.Wrap(err, "listen and serve")
	}

	return nil
}
