package html

import (
	"fmt"
	"html"
	"os"
	"regexp"
	"strings"

	"github.com/verbumby/verbum/backend/ctl/dictimport/dictparser"
	"github.com/verbumby/verbum/backend/textutil"
)

func ParseFile(filename string) (dictparser.Dictionary, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return dictparser.Dictionary{}, fmt.Errorf("read %s: %w", filename, err)
	}
	return ParseString(string(bytes))
}

func ParseString(content string) (dictparser.Dictionary, error) {
	stylesEnd := "</style>\n\n"
	stylesEndPos := strings.LastIndex(content, stylesEnd)
	if stylesEndPos > -1 {
		content = content[stylesEndPos+len(stylesEnd):]
	}

	if content[len(content)-1] == '\n' {
		content = content[:len(content)-1]
	}

	bodies := strings.Split(content, "\n<hr/>\n")

	articles := []dictparser.Article{}
	for _, b := range bodies {
		id, title, hws, hwsalt, err := parseArticleAttributes(b)
		if err != nil {
			return dictparser.Dictionary{}, fmt.Errorf("parse article attributes %s: %w", b, err)
		}
		b = strings.ReplaceAll(b, "<sup>", "\u200b<sup>")
		articles = append(articles, dictparser.Article{
			ID:           id,
			Title:        title,
			Headwords:    hws,
			HeadwordsAlt: hwsalt,
			Body:         b,
		})
	}

	dup := map[string]dictparser.Article{}
	for _, a := range articles {
		if otherA, ok := dup[a.ID]; ok {
			// fmt.Println("dup: ", a.ID)
			return dictparser.Dictionary{}, fmt.Errorf("duplicate id of article %v and %v", otherA, a)
		}
		dup[a.ID] = a
	}

	return dictparser.Dictionary{
		IDsProvided: true,
		Articles:    articles,
	}, nil
}

var reAttr = regexp.MustCompile(`(?m)^<p><strong>(.*?)(?:<sup>(\d+)</sup>)?</strong>`)

func parseArticleAttributes(body string) (id, title string, hws, hwsalt []string, err error) {
	ms := reAttr.FindAllStringSubmatch(body, -1)
	if len(ms) == 0 {
		err = fmt.Errorf("can't find any attributes in %s", body)
		return
	}

	m := ms[0]
	ms = ms[1:]

	hws = parseHeadwords(m[1])

	id = hws[0]
	id = textutil.Slugify(id)
	title = strings.TrimSpace(m[1])

	if m[2] != "" {
		id += "-" + m[2]
		title += " " + m[2]
	}

	for _, m := range ms {
		hwsalt = append(hwsalt, parseHeadwords(m[1])...)
	}

	return
}

var reExpandParentheses = regexp.MustCompile(`\(([^)]*)\)`)

func parseHeadwords(s string) []string {
	s = html.UnescapeString(s)
	hws := strings.Split(s, ",")
	for i := range hws {
		hws[i] = strings.TrimSpace(hws[i])
	}

	expanded := []string{}
	for len(hws) > 0 {
		hw := hws[0]
		hws = hws[1:]

		if !strings.ContainsRune(hw, '(') {
			expanded = append(expanded, hw)
			continue
		}

		hws = append(hws, reExpandParentheses.ReplaceAllString(hw, ""))
		hws = append(hws, reExpandParentheses.ReplaceAllString(hw, "$1"))
	}
	hws = expanded

	return hws
}
