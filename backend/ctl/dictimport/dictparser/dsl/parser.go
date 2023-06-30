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
					a, err := prepareArticle(hws, body.String())
					if err != nil {
						close(articlesCh)
						errCh <- err
						close(errCh)
						return
					}
					articlesCh <- a

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

		a, err := prepareArticle(hws, body.String())
		if err != nil {
			close(articlesCh)
			errCh <- err
			close(errCh)
			return
		}

		articlesCh <- a
		close(articlesCh)
		errCh <- sc.Err()
		close(errCh)
	}()

	return articlesCh, errCh
}

var reTags = regexp.MustCompile(`\[.*?\]`)

func prepareArticle(hwsRaw []string, body string) (dictparser.Article, error) {
	bodyFirstLine := strings.ToLower(firstLine(body))
	bodyFirstLine = reTags.ReplaceAllLiteralString(bodyFirstLine, "")

	var hws []string
	var hwsalt []string

	for _, hw := range prepareHeadwordsForIndexing(hwsRaw) {
		hwLower := strings.ToLower(hw)
		hwInFirstLine := strings.Contains(bodyFirstLine, hwLower)
		if !hwInFirstLine {
			hwsalt = append(hwsalt, hw)
		} else {
			hws = append(hws, hw)
		}
	}

	if len(hws) == 0 {
		hws = hwsalt
		hwsalt = nil
	}

	if len(hws) == 0 {
		return dictparser.Article{}, fmt.Errorf("no headwords for article %v %s found", hwsRaw, body)
	}

	return dictparser.Article{
		Title:        assembleTitleFromHeadwords(hwsRaw),
		Headwords:    hws,
		HeadwordsAlt: hwsalt,
		Phrases:      []string{},
		Body:         body,
	}, nil
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
