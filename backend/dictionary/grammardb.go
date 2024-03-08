package dictionary

import (
	"html/template"
	"regexp"
	"strings"
)

type GrammarDB struct {
	Common
}

func (d GrammarDB) ToHTML(content string) template.HTML {
	re := regexp.MustCompile(`.\x{0301}`)
	substitution := "<v-accent>$0</v-accent>"

	content = re.ReplaceAllString(content, substitution)

	content = strings.ReplaceAll(content, `<table`, `<div class="table-responsive"><table`)
	content = strings.ReplaceAll(content, `</table>`, `</table></div>`)

	if d.abbrevs != nil {
		content = renderAbbrevs(content, d.abbrevs)
	}

	return template.HTML(content)
}
