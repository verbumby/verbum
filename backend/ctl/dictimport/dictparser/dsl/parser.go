package dsl

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/verbumby/verbum/backend/ctl/dictimport/dictparser"
)

func ParseDSLFile(filename string) (dictparser.Dictionary, error) {
	f, err := os.Open(filename)
	if err != nil {
		return dictparser.Dictionary{}, fmt.Errorf("open %s file: %w", filename, err)
	}
	defer f.Close()

	return ParseDSLReader(filename, f)
}

func ParseDSLReader(filename string, r io.Reader) (dictparser.Dictionary, error) {
	ditf, err := ParseReader(filename, r)
	if err != nil {
		return dictparser.Dictionary{}, err
	}

	indentRe := regexp.MustCompile(`(?m)^\t`)

	d := ditf.(dictparser.Dictionary)

	for i, a := range d.Articles {
		a.Body = indentRe.ReplaceAllLiteralString(a.Body, "")
		bodyLower := strings.ToLower(a.Body)
		bodyFirstLine := firstLine(bodyLower)

		hws := []string{}
		hwsalt := []string{}

		for _, hw := range prepareHeadwordsForIndexing(a.Headwords) {
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

		if len(hws) == 0 {
			return d, fmt.Errorf("no headwords for article %v found", a)
		}

		d.Articles[i].Title = assembleTitleFromHeadwords(a.Headwords)
		d.Articles[i].Headwords = hws
		d.Articles[i].HeadwordsAlt = hwsalt
		d.Articles[i].Phrases = []string{}
		d.Articles[i].Body = a.Body
	}

	return d, nil
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
