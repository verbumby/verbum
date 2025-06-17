package dictionary

import (
	"html/template"
	"strings"
)

type HTML struct {
	Common
}

func (d HTML) ToHTML(content string) template.HTML {
	content = wrapAccentedChars(content)

	content = renderAbbrevs(content, d.abbrevs)

	content = strings.ReplaceAll(content, `href="#`, `target="_blank" href="/`+d.ID()+`/`)

	return template.HTML(content)
}
