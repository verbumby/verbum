package main

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
)

func main() {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(errors.Wrap(err, "listen and serve"))
	}
}
