package dictionary

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/verbumby/verbum/backend/config"
	"golang.org/x/text/unicode/norm"
)

type Abbrevs struct {
	list  []*Abbrev
	cache map[string]*Abbrev
}

func (a *Abbrevs) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.list)
}

type Abbrev struct {
	Keys  []string
	Value string
}

func loadDSLAbbrevs(filename string) (*Abbrevs, error) {
	filename = config.DictsRepoPath() + "/" + filename
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", filename, err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

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
			if c.Value != "" {
				c.Value += "\n"
			}
			c.Value += line
		} else {
			if keysSealed {
				list = append(list, c)
				c = &Abbrev{}
				keysSealed = false
			}

			key := norm.NFC.String(strings.TrimSpace(line))
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
			if _, ok := cache[k]; !ok {
				cache[k] = c
			}
		}
	}

	return &Abbrevs{
		list:  list,
		cache: cache,
	}, nil
}

var (
	reAbbrev    = regexp.MustCompile(`(?U)<v-abbr.*</v-abbr>`)
	reStripHtml = regexp.MustCompile(`(?U)</?.*>`)
)

func renderAbbrevs(content string, abbrevs *Abbrevs) string {
	if abbrevs == nil {
		return content
	}

	return reAbbrev.ReplaceAllStringFunc(content, func(m string) string {
		text := reStripHtml.ReplaceAllLiteralString(m, "")
		text = norm.NFC.String(text)
		if v, ok := abbrevs.cache[text]; ok {
			tt := `<v-abbr data-bs-toggle="tooltip" data-bs-title="%s" tabindex="0"`
			m = strings.Replace(m, "<v-abbr", fmt.Sprintf(tt, html.EscapeString(v.Value)), 1)
		}
		return m
	})
}
