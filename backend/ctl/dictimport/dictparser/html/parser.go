package html

import (
	"bufio"
	"fmt"
	"html"
	"io"
	"regexp"
	"strings"

	"github.com/verbumby/verbum/backend/ctl/dictimport/dictparser"
	"github.com/verbumby/verbum/backend/dictionary"
	"github.com/verbumby/verbum/backend/textutil"
	"golang.org/x/text/unicode/norm"
)

func ParseReader(r io.Reader, settings dictionary.IndexSettings) (chan dictparser.Article, chan error) {
	articlesCh := make(chan dictparser.Article, 64)
	errCh := make(chan error)

	go func() {
		sc := bufio.NewScanner(r)
		sc.Buffer(make([]byte, 16*1024), bufio.MaxScanTokenSize*4)
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

			a, err := parseArticle(bodyStr, settings)
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
	reHW    = regexp.MustCompile(`<(?:strong|b) class="(hw(?:-alt)?)"(?: id="([^"]+)")?[^>]*>([^<]*)</(?:strong|b)>`)
	reIndex = regexp.MustCompile(`<sup[^>]*>([\dIVX-]+)</sup>`)
)

func parseArticle(body string, settings dictionary.IndexSettings) (dictparser.Article, error) {
	ms := reHW.FindAllStringSubmatch(body, -1)
	if len(ms) == 0 {
		return dictparser.Article{}, fmt.Errorf("can't find any attributes in %s", body)
	}

	hws := []string{}
	idAttr := ""
	var hwsalt []string

	for _, m := range ms {
		hw := m[3]
		hw = strings.TrimSpace(hw)
		hw = norm.NFD.String(hw)
		hw = strings.ReplaceAll(hw, "\u0301", "")
		hw = strings.ReplaceAll(hw, "\u0311", "")
		hw = strings.ReplaceAll(hw, "\u030c", "")
		hw = norm.NFC.String(hw)

		parenthesisBalance := 0
		errParenthesisUnbalanced := fmt.Errorf("headword `%s` has unbalanced or out of order parenthesis", hw)
		for _, r := range hw {
			switch r {
			case '(':
				parenthesisBalance++
			case ')':
				parenthesisBalance--
			}
			if parenthesisBalance < 0 {
				return dictparser.Article{}, errParenthesisUnbalanced
			}
		}
		if parenthesisBalance != 0 {
			return dictparser.Article{}, errParenthesisUnbalanced
		}

		switch m[1] {
		case "hw":
			hws = append(hws, hw)
		case "hw-alt":
			hwsalt = append(hwsalt, hw)
		}

		if m[2] != "" {
			idAttr = m[2]
		}
	}

	if len(hws) == 0 {
		return dictparser.Article{}, fmt.Errorf("no headwords found in %s", body)
	}

	title := strings.Join(hws, ", ")

	hws = expandHeadwords(hws)
	hwsalt = expandHeadwords(hwsalt)

	for i := range hws {
		hws[i] = html.UnescapeString(hws[i])
	}
	for i := range hwsalt {
		hwsalt[i] = html.UnescapeString(hwsalt[i])
	}

	id := idAttr
	if id == "" {
		id = hws[0]

		idx := ""
		if m := reIndex.FindStringSubmatch(body); m != nil {
			idx = m[1]
		}

		if idx != "" {
			id += "-" + idx
			title += " " + idx
		}
	}


	body = strings.ReplaceAll(body, "<sup", "\u200b<sup")
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
