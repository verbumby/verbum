package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/verbumby/verbum/backend/pkg/article"
	"github.com/verbumby/verbum/backend/pkg/db"
	reform "gopkg.in/reform.v1"
)

var rvblrNumlistsFixCmd = &cobra.Command{
	Use: "rvblr-numlists-fix",
	Run: func(cmd *cobra.Command, args []string) {
		tail := "INNER JOIN tasks_articles_rel " +
			"ON articles.id = tasks_articles_rel.article_id " +
			"AND tasks_articles_rel.status = 'PENDING' " +
			"AND tasks_articles_rel.task_id = 1 " //+
			//"LIMIT 100 OFFSET 1234"
		records, err := db.DB.SelectAllFrom(article.ArticleTable, tail)
		if err != nil {
			log.Fatal(err)
		}
		var re = regexp.MustCompile(`\s+(\d+)\.\s+`)
		for i, record := range records {
			a := record.(*article.Article)
			c := re.ReplaceAllString(a.Content, "\n\n$1. ")
			c = re.ReplaceAllString(c, "\n\n$1. ")
			a.Content = c
			if err := db.DB.Update(a); err != nil && err != reform.ErrNoRows {
				log.Fatal(err)
			}
			fmt.Print(".")
			if i%80 == 0 {
				fmt.Println()
			}
		}
	},
}
