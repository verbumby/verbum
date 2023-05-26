package stardict

import (
	"fmt"
	"os"
	"strings"

	"github.com/verbumby/verbum/backend/pkg/ctl/dictimport/dictparser"
)

func LoadArticles(filename string) (dictparser.Dictionary, error) {
	contentBytes, err := os.ReadFile(filename)
	if err != nil {
		return dictparser.Dictionary{}, fmt.Errorf("read file %s: %w", filename, err)
	}

	content := string(contentBytes)
	articles := []dictparser.Article{}

	articleSources := strings.Split(content, "<k>")
	for _, articleSource := range articleSources {
		if strings.TrimSpace(articleSource) == "" {
			continue
		}

		title, body, found := strings.Cut(articleSource, "</k>")
		if !found {
			return dictparser.Dictionary{}, fmt.Errorf("expected a key-body separator </k> in article: %s", articleSource)
		}

		hws := strings.Split(title, ",")
		for i, hw := range hws {
			hws[i] = strings.TrimSpace(hw)
		}

		articles = append(articles, dictparser.Article{
			Title:     title,
			Headwords: hws,
			Body:      "<k>" + title + "</k>\n" + body,
		})
	}

	return dictparser.Dictionary{
		Articles: articles,
	}, nil
}
