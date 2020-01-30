package dsl

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/verbumby/verbum/backend/pkg/ctl/dictimport/dictparser"
)

func ParseDSLReader(filename string, r io.Reader) (dictparser.Dictionary, error) {
	ditf, err := ParseReader(filename, r)
	if err != nil {
		return dictparser.Dictionary{}, err
	}

	var indentRe = regexp.MustCompile(`(?m)^\t`)
	var hwCurlyBracesRe = regexp.MustCompile(`(?U){.*}`)

	d := ditf.(dictparser.Dictionary)

	for i, a := range d.Articles {
		a.Body = indentRe.ReplaceAllLiteralString(a.Body, "")
		bodylower := strings.ToLower(a.Body)

		phws := []string{}
		for _, phw := range a.Headwords {
			phw = hwCurlyBracesRe.ReplaceAllLiteralString(phw, "")
			phws = append(phws, phw)
		}
		a.Headwords = phws

		hws := []string{}
		hwsalt := []string{}

		for _, phw := range a.Headwords {
			phw = strings.TrimSpace(phw)
			phw = strings.ReplaceAll(phw, "\\(", "(")
			phw = strings.ReplaceAll(phw, "\\)", ")")

			phrase := fmt.Sprintf("[b][ex][lang id=1049][c steelblue]%s[/c][/lang][/ex][/b]", phw)
			phrase = strings.ToLower(phrase)

			phw = strings.ReplaceAll(phw, "...", "")

			if strings.Contains(bodylower, phrase) || strings.Contains(phw, "(") {
				hwsalt = append(hwsalt, phw)
			} else {
				hws = append(hws, phw)
			}
		}

		d.Articles[i].Headwords = hws
		d.Articles[i].HeadwordsAlt = hwsalt
		d.Articles[i].Body = a.Body
	}

	return d, nil
}
