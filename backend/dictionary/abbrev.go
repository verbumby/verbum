package dictionary

import (
	"bufio"
	_ "embed"
	"fmt"
	"html"
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

func parseDSLAbbrev(content string) map[string]string {
	s := bufio.NewScanner(strings.NewReader(content))
	result := map[string]string{}

	keys := []string{}
	keysSealed := false

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
			for _, key := range keys {
				if _, ok := result[key]; ok {
					line = "\n" + line
				}
				result[key] += line
			}
		} else {
			if keysSealed {
				keys = []string{}
				keysSealed = false
			}

			key := strings.TrimSpace(line)
			if _, ok := result[key]; ok {
				fmt.Println("duplicate abbrev key: " + key)
			}

			keys = append(keys, key)
		}
	}

	if s.Err() != nil {
		panic(s.Err())
	}

	extra := map[string]string{}
	for k, v := range result {
		kr := []rune(k)
		kr[0] = unicode.ToUpper(kr[0])
		k = string(kr)

		if _, ok := result[k]; !ok {
			extra[k] = v
		}
	}

	for k, v := range extra {
		result[k] = v
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
		if v, ok := abbrevs[text]; ok {
			tt := `<v-abbr data-bs-toggle="tooltip" data-bs-title="%s" tabindex="0">`
			m = strings.Replace(m, "<v-abbr>", fmt.Sprintf(tt, html.EscapeString(v)), 1)
		}
		return m
	})
}
