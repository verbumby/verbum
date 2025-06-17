package dictionary

import (
	"html/template"
	"strings"
)

type GrammarDB struct {
	Common
}

func (d GrammarDB) ToHTML(content string) template.HTML {
	content = wrapAccentedChars(content)

	content = strings.ReplaceAll(content, `<table`, `<div class="table-responsive"><table`)
	content = strings.ReplaceAll(content, `</table>`, `</table></div>`)

	if d.abbrevs != nil {
		content = renderAbbrevs(content, d.abbrevs)
	}

	return template.HTML(content)
}
