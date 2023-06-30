package stardict

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/verbumby/verbum/backend/ctl/dictimport/dictparser"
)

func LoadArticles(r io.Reader) (chan dictparser.Article, chan error) {
	articlesCh := make(chan dictparser.Article, 64)
	errCh := make(chan error)

	go func() {
		sc := bufio.NewScanner(r)
		sc.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			if atEOF && len(data) == 0 {
				return 0, nil, nil
			}

			delim := []byte("<k>")

			if i := bytes.Index(data, delim); i >= 0 {
				return i + len(delim), data[0:i], nil
			}

			if atEOF {
				return len(data), data, nil
			}

			return 0, nil, nil
		})

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
