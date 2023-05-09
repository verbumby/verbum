package dictionary

import (
	"bufio"
	_ "embed"
	"fmt"
	"html"
	"regexp"
	"strings"
)

//go:embed esbm_abrv.dsl
var esbmAbbrev string

//go:embed tsbm_abrv.dsl
var tsbmAbbrev string

//go:embed brs_abrv.dsl
var brsAbbrev string

//go:embed rbs_abrv.dsl
var rbsAbbrev string

func parseDSLAbbrev(content string) map[string]string {
	s := bufio.NewScanner(strings.NewReader(content))
	result := map[string]string{}
	key := ""
	for s.Scan() {
		line := s.Text()

		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.TrimSpace(line) == "" {
			continue
		}

		if strings.HasPrefix(line, "\t") || strings.HasPrefix(line, " ") {
			line = strings.TrimSpace(line)
			if _, ok := result[key]; ok {
				line = "\n" + line
			}
			result[key] += strings.TrimSpace(line)
		} else {
			key = strings.ToLower(strings.TrimSpace(line))
			if _, ok := result[key]; ok {
				fmt.Println("duplicate abbrev key: " + key)
				// panic("duplicate abbrev key: " + key)
			}
		}
	}

	if s.Err() != nil {
		panic(s.Err())
	}

	return result
}

var (
	reAbbrev    = regexp.MustCompile(`(?U)<v-abbr>.*</v-abbr>`)
	reStripHtml = regexp.MustCompile(`(?U)</?.*>`)
)

func renderAbbrevs(content string, abbrevs map[string]string) string {
	return reAbbrev.ReplaceAllStringFunc(content, func(m string) string {
		text := reStripHtml.ReplaceAllLiteralString(m, "")
		text = strings.ToLower(text)
		if v, ok := abbrevs[text]; ok {
			m = strings.Replace(m, "<v-abbr>", fmt.Sprintf(`<v-abbr title="%s">`, html.EscapeString(v)), 1)
		}
		return m
	})
}
