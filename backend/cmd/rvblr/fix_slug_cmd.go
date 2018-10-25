package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/verbumby/verbum/backend/pkg/db"
	"github.com/verbumby/verbum/backend/pkg/dict"
	"github.com/verbumby/verbum/backend/pkg/utils"
)

var fixSlugCmd = &cobra.Command{
	Use: "fix-slug",
	Run: func(cmd *cobra.Command, args []string) {
		structs, err := db.DB.SelectAllFrom(dict.DictTable, "")
		if err != nil {
			log.Fatalf("select all dicts: %v", err)
		}
		for _, s := range structs {
			d := s.(*dict.Dict)
			slug := utils.Slugify(d.Title)
			if slug != d.Slug {
				fmt.Println(d.Title, "-->", slug)
				d.Slug = slug
				if fixSlugCommit {
					if err := db.DB.Save(d); err != nil {
						log.Fatalf("save dict %d: %v", d.ID, err)
					}
				}
			}
		}
	},
}
var fixSlugCommit = false

func init() {
	fixSlugCmd.Flags().BoolVar(&fixSlugCommit, "commit", false, "whether to commit changes to the database")
}
