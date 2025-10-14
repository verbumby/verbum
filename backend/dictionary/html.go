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

	content = strings.ReplaceAll(content, `<img src="`, `<img src="/images/`+d.ID()+`/`)

	content = strings.ReplaceAll(content, `<a href="img/`, `<a href="/images/`+d.ID()+`/img/`)

	return template.HTML(content)
}
