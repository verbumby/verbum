package html

import (
	"fmt"
	"html"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/verbumby/verbum/backend/ctl/dictimport/dictparser"
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
	dups := false
	for _, a := range articles {
		if otherA, ok := dup[a.ID]; ok {
			dups = true
			log.Printf("duplicate id of article %v and %v", otherA, a)
		}
		dup[a.ID] = a
	}

	if dups {
		return dictparser.Dictionary{}, fmt.Errorf("there are duplicate ids")
	}

	return dictparser.Dictionary{
		IDsProvided: true,
		Articles:    articles,
	}, nil
}

var (
	reHW    = regexp.MustCompile(`<strong class="(hw(?:-alt)?)">([^<]*)</strong>`)
	reIndex = regexp.MustCompile(`<sup>(\d+)</sup>`)
)

func parseArticleAttributes(body string) (id, title string, hws, hwsalt []string, err error) {
	ms := reHW.FindAllStringSubmatch(body, -1)
	if len(ms) == 0 {
		err = fmt.Errorf("can't find any attributes in %s", body)
		return
	}

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
		err = fmt.Errorf("no headwords found in %s", body)
		return
	}

	title = strings.Join(hws, ", ")

	hws = expandHeadwords(hws)
	hwsalt = expandHeadwords(hwsalt)

	id = hws[0]

	idx := ""
	if m := reIndex.FindStringSubmatch(body); m != nil {
		idx = m[1]
	}

	if idx != "" {
		id += "-" + idx
		title += " " + idx
	}

	return
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
