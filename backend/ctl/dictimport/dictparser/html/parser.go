package html

import (
	"bufio"
	"fmt"
	"html"
	"io"
	"regexp"
	"strings"

	"github.com/verbumby/verbum/backend/ctl/dictimport/dictparser"
	"github.com/verbumby/verbum/backend/textutil"
)

func ParseReader(r io.Reader) (chan dictparser.Article, chan error) {
	articlesCh := make(chan dictparser.Article, 64)
	errCh := make(chan error)

	go func() {
		sc := bufio.NewScanner(r)
		sc.Split(textutil.GetDelimSplitFunc("<hr/>\n"))
		firstArticle := true

		for sc.Scan() {
			bodyStr := sc.Text()

			if firstArticle {
				firstArticle = false

				stylesEnd := "</style>\n\n"
				stylesEndPos := strings.LastIndex(bodyStr, stylesEnd)
				if stylesEndPos > -1 {
					bodyStr = bodyStr[stylesEndPos+len(stylesEnd):]
				}
			}

			a, err := parseArticle(bodyStr)
			if err != nil {
				close(articlesCh)
				errCh <- fmt.Errorf("parse article %s: %w", bodyStr, err)
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

var (
	reHW    = regexp.MustCompile(`<strong class="(hw(?:-alt)?)">([^<]*)</strong>`)
	reIndex = regexp.MustCompile(`<sup>(\d+)</sup>`)
)

func parseArticle(body string) (dictparser.Article, error) {
	ms := reHW.FindAllStringSubmatch(body, -1)
	if len(ms) == 0 {
		return dictparser.Article{}, fmt.Errorf("can't find any attributes in %s", body)
	}

	hws := []string{}
	var hwsalt []string

	for _, m := range ms {
		hw := m[2]
		hw = html.UnescapeString(hw)
		hw = strings.TrimSpace(hw)

		switch m[1] {
		case "hw":
			hws = append(hws, hw)
		case "hw-alt":
			hwsalt = append(hwsalt, hw)
		}
	}

	if len(hws) == 0 {
		return dictparser.Article{}, fmt.Errorf("no headwords found in %s", body)
	}

	title := strings.Join(hws, ", ")

	hws = expandHeadwords(hws)
	hwsalt = expandHeadwords(hwsalt)

	id := hws[0]

	idx := ""
	if m := reIndex.FindStringSubmatch(body); m != nil {
		idx = m[1]
	}

	if idx != "" {
		id += "-" + idx
		title += " " + idx
	}

	body = strings.ReplaceAll(body, "<sup>", "\u200b<sup>")
	return dictparser.Article{
		ID:           id,
		Title:        title,
		Headwords:    hws,
		HeadwordsAlt: hwsalt,
		Body:         body,
	}, nil
}

var reExpandParentheses = regexp.MustCompile(`\(([^)]*)\)`)

func expandHeadwords(hws []string) []string {
	if hws == nil {
		return nil
	}

	expanded := []string{}
	for len(hws) > 0 {
		hw := hws[0]
		hws = hws[1:]

		if !strings.ContainsRune(hw, '(') {
			expanded = append(expanded, hw)
			continue
		}

		hws = append(hws, strings.TrimSpace(reExpandParentheses.ReplaceAllString(hw, "")))
		hws = append(hws, strings.TrimSpace(reExpandParentheses.ReplaceAllString(hw, "$1")))
	}

	return expanded
}
