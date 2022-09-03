package dsl

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/verbumby/verbum/backend/pkg/ctl/dictimport/dictparser"
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

	var indentRe = regexp.MustCompile(`(?m)^\t`)

	d := ditf.(dictparser.Dictionary)

	for i, a := range d.Articles {
		a.Body = indentRe.ReplaceAllLiteralString(a.Body, "")

		hws := []string{}
		hwsalt := []string{}

		for _, phw := range a.Headwords {
			phw = strings.TrimSpace(phw)
			phw = strings.ReplaceAll(phw, "\\(", "(")
			phw = strings.ReplaceAll(phw, "\\)", ")")
			phw = strings.ReplaceAll(phw, "...", "")

			hws = append(hws, phw)
		}

		title := strings.Join(hws, ", ")
		title = strings.ReplaceAll(title, "{", "")
		title = strings.ReplaceAll(title, "}", "")
		title = strings.ReplaceAll(title, "\\~", "~")
		title = strings.ReplaceAll(title, " ,", ",")

		var reCurlyBrace = regexp.MustCompile(`{.*?}`)
		d.Articles[i].Headwords = []string{}
		for _, hw := range hws {
			hw = reCurlyBrace.ReplaceAllString(hw, "")
			d.Articles[i].Headwords = append(d.Articles[i].Headwords, hw)
		}

		if len(hws) == 0 {
			return d, fmt.Errorf("no headwords for article %v found", a)
		}

		d.Articles[i].Title = title
		d.Articles[i].HeadwordsAlt = hwsalt
		d.Articles[i].Phrases = []string{}
		d.Articles[i].Body = a.Body
	}

	return d, nil
}
