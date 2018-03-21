package main

import (
	"log"

	"github.com/verbumby/verbum/backend/pkg/app"
)

func main() {
	if err := app.Bootstrap(); err != nil {
		log.Fatal(err)
	}
}
