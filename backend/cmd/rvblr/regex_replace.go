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

// RegexReplaceCmdFlags this command flags
var RegexReplaceCmdFlags = struct {
	Regex       string
	Replacement string
	Persist     bool
}{}

var regexReplaceCmd = &cobra.Command{
	Use: "regex-replace",
	Run: func(cmd *cobra.Command, args []string) {
		tail := "INNER JOIN tasks_articles_rel " +
			"ON articles.id = tasks_articles_rel.article_id " +
			"AND tasks_articles_rel.status = 'PENDING' " +
			"AND tasks_articles_rel.task_id = 1 "
		if !RegexReplaceCmdFlags.Persist {
			tail += "LIMIT 100 OFFSET 1234 "
		}
		records, err := db.DB.SelectAllFrom(article.ArticleTable, tail)
		if err != nil {
			log.Fatal(err)
		}
		var re = regexp.MustCompile(RegexReplaceCmdFlags.Regex)
		for i, record := range records {
			a := record.(*article.Article)
			c := re.ReplaceAllString(a.Content, RegexReplaceCmdFlags.Replacement)
			c = re.ReplaceAllString(c, RegexReplaceCmdFlags.Replacement)
			a.Content = c

			if RegexReplaceCmdFlags.Persist {
				if err := db.DB.Update(a); err != nil && err != reform.ErrNoRows {
					log.Fatal(err)
				}
				fmt.Print(".")
				if i%80 == 0 {
					fmt.Println()
				}
			} else {
				fmt.Println("-----------------------------------")
				fmt.Println(a.Content)
			}
		}
	},
}

func init() {
	regexReplaceCmd.Flags().StringVarP(&RegexReplaceCmdFlags.Regex, "regex", "e", "", "Regular Expression")
	regexReplaceCmd.Flags().StringVarP(&RegexReplaceCmdFlags.Replacement, "replacement", "r", "", "Replacement")
	regexReplaceCmd.Flags().BoolVarP(&RegexReplaceCmdFlags.Persist, "persist", "p", false, "Whether to persist changes")
}
