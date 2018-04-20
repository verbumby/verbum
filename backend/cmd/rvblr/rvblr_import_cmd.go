package main

import (
	"encoding/xml"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/verbumby/verbum/backend/pkg/article"
	"github.com/verbumby/verbum/backend/pkg/db"
	"github.com/verbumby/verbum/backend/pkg/dict"
)

var rvblrImportCmd = &cobra.Command{
	Use: "rvblr-import",
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
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
	},
}
