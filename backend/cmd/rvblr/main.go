package main

import (
	"encoding/xml"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/verbumby/verbum/backend/pkg/article"
	"github.com/verbumby/verbum/backend/pkg/db"
	"github.com/verbumby/verbum/backend/pkg/dict"

	"github.com/pkg/errors"
	"github.com/verbumby/verbum/backend/pkg/app"
)

type DictData struct {
	Items ItemsData `xml:"items"`
}

type ItemsData struct {
	Item []ItemData `xml:"item"`
}

type ItemData struct {
	Title string `xml:"title"`
	Meta  string `xml:"meta"`
	Def   string `xml:"definition"`
}

func main() {
	var filePath string
	flag.StringVar(&filePath, "file", "", "path to the rv-blr.xml file")
	flag.Parse()

	if err := app.Bootstrap(); err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(errors.Wrap(err, "open file"))
	}

	var dd DictData
	if err := xml.NewDecoder(f).Decode(&dd); err != nil {
		log.Fatal(errors.Wrap(err, "decode xml"))
	}

	d := &dict.Dict{
		Title: "Тлумачальны слоўнік беларускай мовы (rv-blr.com)",
	}
	if err := db.DB.Save(d); err != nil {
		log.Fatal(errors.Wrap(err, "create dictionary"))
	}

	for _, item := range dd.Items.Item {
		article := &article.Article{
			Title:   strings.ToLower(item.Title),
			Content: "<v-hw>" + strings.ToLower(item.Title) + "</v-hw>\n\n*" + item.Meta + "*\n\n" + item.Def,
			DictID:  d.ID,
		}
		db.DB.Save(article)
	}
}
