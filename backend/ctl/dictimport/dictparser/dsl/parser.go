package dsl

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/verbumby/verbum/backend/ctl/dictimport/dictparser"
)

func ParseReader(r io.Reader) (chan dictparser.Article, chan error) {
	sc := bufio.NewScanner(r)

	hws := []string{}

	for sc.Scan() {
		t := sc.Text()
		if strings.TrimSpace(t) == "" {
			continue
		}

		if t[0] == '#' {
			continue
		}

		hws = append(hws, t)
		break
	}

	articlesCh := make(chan dictparser.Article, 64)
	errCh := make(chan error)

	go func() {
		hwsSealed := false
		body := strings.Builder{}

		for sc.Scan() {
			t := sc.Text()

			if len(t) > 0 && t[0] != '\t' {
				if hwsSealed {
					articlesCh <- prepareArticle(hws, body.String())
					hws = hws[0:0]
					body.Reset()
					hwsSealed = false
				}
				hws = append(hws, t)
			} else {
				hwsSealed = true
				if len(t) > 0 {
					body.WriteString(t[1:])
				}
				body.WriteRune('\n')
			}
		}

		articlesCh <- prepareArticle(hws, body.String())
		close(articlesCh)
		errCh <- sc.Err()
		close(errCh)
	}()

	return articlesCh, errCh
}

func prepareArticle(hwsRaw []string, body string) dictparser.Article {
	bodyLower := strings.ToLower(body)
	bodyFirstLine := firstLine(bodyLower)

	hws := []string{}
	hwsalt := []string{}

	for _, hw := range prepareHeadwordsForIndexing(hwsRaw) {
		hwLower := strings.ToLower(hw)
		ex := fmt.Sprintf("[ex][lang id=1049][c steelblue]%s[/c][/lang][/ex]", hwLower)
		exContains := strings.Contains(bodyLower, ex)
		hwInFirstLine := strings.Contains(bodyFirstLine, hwLower)
		if exContains && !hwInFirstLine {
			hwsalt = append(hwsalt, hw)
		} else {
			hws = append(hws, hw)
		}
	}

	// if len(hws) == 0 {
	// 	return d, fmt.Errorf("no headwords for article %v found", a)
	// }

	return dictparser.Article{
		Title:        assembleTitleFromHeadwords(hwsRaw),
		Headwords:    hws,
		HeadwordsAlt: hwsalt,
		Phrases:      []string{},
		Body:         body,
	}
}

func firstLine(s string) string {
	nl := strings.IndexRune(s, '\n')
	if nl == -1 {
		nl = len(s)
	}

	return s[:nl]
}

var reCurlyBrace = regexp.MustCompile(`{.*?}`)

func prepareHeadwordsForIndexing(hws []string) []string {
	result := []string{}
	for _, hw := range hws {
		hw = reCurlyBrace.ReplaceAllString(hw, "")
		hw = strings.ReplaceAll(hw, "\\(", "(")
		hw = strings.ReplaceAll(hw, "\\)", ")")
		hw = strings.ReplaceAll(hw, "...", "")
		hw = strings.TrimSpace(hw)
		result = append(result, hw)
	}

	return result
}

func assembleTitleFromHeadwords(hws []string) string {
	result := []string{}
	for _, hw := range hws {
		hw = strings.TrimSpace(hw)
		hw = strings.ReplaceAll(hw, "{", "")
		hw = strings.ReplaceAll(hw, "}", "")
		hw = strings.ReplaceAll(hw, "\\(", "(")
		hw = strings.ReplaceAll(hw, "\\)", ")")
		hw = strings.ReplaceAll(hw, "\\~", "~")
		hw = strings.ReplaceAll(hw, " ,", ",")
		result = append(result, hw)
	}

	nodup := []string{}
outer:
	for _, hw := range result {
		for _, noduphw := range nodup {
			if hw == noduphw {
				continue outer
			}
		}
		nodup = append(nodup, hw)
	}

	nodup = nodup[:min(3, len(nodup))]
	return strings.Join(nodup, ", ")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
