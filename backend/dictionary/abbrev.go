package dictionary

import (
	"bufio"
	_ "embed"
	"fmt"
	"html"
	"log"
	"regexp"
	"strings"
	"unicode"
)

//go:embed esbm_abrv.dsl
var esbmAbbrev string

//go:embed tsbm_abrv.dsl
var tsbmAbbrev string

//go:embed brs_abrv.dsl
var brsAbbrev string

//go:embed rbs_abrv.dsl
var rbsAbbrev string

type Abbrevs struct {
	list  []*Abbrev
	cache map[string]*Abbrev
}

type Abbrev struct {
	Keys  []string
	Value string
}

func loadDSLAbbrevs(content string) (*Abbrevs, error) {
	s := bufio.NewScanner(strings.NewReader(content))

	keysSealed := false
	list := []*Abbrev{}
	c := &Abbrev{}

	for s.Scan() {
		line := s.Text()

		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.TrimSpace(line) == "" {
			continue
		}

		if strings.HasPrefix(line, "\t") || strings.HasPrefix(line, " ") {
			keysSealed = true
			line = strings.TrimSpace(line)
			c.Value += line + "\n"
		} else {
			if keysSealed {
				list = append(list, c)
				c = &Abbrev{}
				keysSealed = false
			}

			key := strings.TrimSpace(line)
			c.Keys = append(c.Keys, key)
		}
	}

	list = append(list, c)

	if s.Err() != nil {
		return nil, fmt.Errorf("scanner error: %w", s.Err())
	}

	cache := map[string]*Abbrev{}

	for _, c := range list {
		for _, k := range c.Keys {
			if _, ok := cache[k]; ok {
				log.Printf("duplicate abbrev key: %s", k)
			}
			cache[k] = c
		}
	}

	for _, c := range list {
		for _, k := range c.Keys {
			kr := []rune(k)
			if unicode.IsUpper(kr[0]) {
				continue
			}
			kr[0] = unicode.ToUpper(kr[0])
			k = string(kr)
			cache[k] = c
		}
	}

	return &Abbrevs{
		list:  list,
		cache: cache,
	}, nil
}

var (
	reAbbrev    = regexp.MustCompile(`(?U)<v-abbr>.*</v-abbr>`)
	reStripHtml = regexp.MustCompile(`(?U)</?.*>`)
)

func renderAbbrevs(content string, abbrevs *Abbrevs) string {
	return reAbbrev.ReplaceAllStringFunc(content, func(m string) string {
		text := reStripHtml.ReplaceAllLiteralString(m, "")
		if v, ok := abbrevs.cache[text]; ok {
			tt := `<v-abbr data-bs-toggle="tooltip" data-bs-title="%s" tabindex="0">`
			m = strings.Replace(m, "<v-abbr>", fmt.Sprintf(tt, html.EscapeString(v.Value)), 1)
		}
		return m
	})
}
