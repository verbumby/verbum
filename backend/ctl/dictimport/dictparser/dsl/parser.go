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
		if len(a.Headwords) == 0 {
			return d, fmt.Errorf("no headwords for article %v found", a)
		}

		a.Body = indentRe.ReplaceAllLiteralString(a.Body, "")

		d.Articles[i].Title = assembleTitleFromHeadwords(a.Headwords)
		d.Articles[i].Headwords = prepareHeadwordsForIndexing(a.Headwords)
		d.Articles[i].HeadwordsAlt = []string{}
		d.Articles[i].Phrases = []string{}
		d.Articles[i].Body = a.Body
	}

	return d, nil
}

var reCurlyBrace = regexp.MustCompile(`{.*?}`)

func prepareHeadwordsForIndexing(hws []string) []string {
	for i := range hws {
		hws[i] = reCurlyBrace.ReplaceAllString(hws[i], "")
		hws[i] = strings.ReplaceAll(hws[i], "\\(", "(")
		hws[i] = strings.ReplaceAll(hws[i], "\\)", ")")
		hws[i] = strings.ReplaceAll(hws[i], "...", "")
		hws[i] = strings.TrimSpace(hws[i])
	}

	return hws
}

func assembleTitleFromHeadwords(hws []string) string {
	for i := range hws {
		hws[i] = strings.TrimSpace(hws[i])
		hws[i] = strings.ReplaceAll(hws[i], "{", "")
		hws[i] = strings.ReplaceAll(hws[i], "}", "")
		hws[i] = strings.ReplaceAll(hws[i], "\\(", "(")
		hws[i] = strings.ReplaceAll(hws[i], "\\)", ")")
		hws[i] = strings.ReplaceAll(hws[i], "\\~", "~")
		hws[i] = strings.ReplaceAll(hws[i], " ,", ",")
	}

	nodup := []string{}
outer:
	for _, hw := range hws {
		for _, noduphw := range nodup {
			if hw == noduphw {
				continue outer
			}
		}
		nodup = append(nodup, hw)
	}
	return strings.Join(nodup, ", ")
}
