package stardict

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/verbumby/verbum/backend/ctl/dictimport/dictparser"
	"github.com/verbumby/verbum/backend/textutil"
)

func LoadArticles(r io.Reader) (chan dictparser.Article, chan error) {
	articlesCh := make(chan dictparser.Article, 64)
	errCh := make(chan error)

	go func() {
		sc := bufio.NewScanner(r)
		sc.Split(textutil.GetDelimSplitFunc("<k>"))
		for sc.Scan() {
			articleSource := sc.Text()
			if strings.TrimSpace(articleSource) == "" {
				continue
			}

			a, err := parseArticle(articleSource)
			if err != nil {
				close(articlesCh)
				errCh <- err
				close(errCh)
				return
			}
			articlesCh <- a
		}

		close(articlesCh)
		close(errCh)
	}()

	return articlesCh, errCh
}

func parseArticle(body string) (dictparser.Article, error) {
	title, body, found := strings.Cut(body, "</k>")
	if !found {
		return dictparser.Article{}, fmt.Errorf("expected a key-body separator </k> in article: %s", body)
	}

	hws := strings.Split(title, ",")
	for i, hw := range hws {
		hws[i] = strings.TrimSpace(hw)
	}

	return dictparser.Article{
		Title:     title,
		Headwords: hws,
		Body:      "<k>" + title + "</k>\n" + strings.TrimSpace(body),
	}, nil
}
