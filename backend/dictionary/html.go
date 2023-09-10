package dictionary

import (
	"html/template"
	"regexp"
	"strings"

	"golang.org/x/text/unicode/norm"
)

type HTML struct {
	Common
}

func (d HTML) ToHTML(content string) template.HTML {
	content = norm.NFD.String(content)

	re := regexp.MustCompile(`.\x{0301}`)
	substitution := "<v-accent>$0</v-accent>"

	content = re.ReplaceAllString(content, substitution)

	content = renderAbbrevs(content, d.abbrevs)

	content = norm.NFC.String(content)

	content = strings.ReplaceAll(content, `href="#`, `target="_blank" href="/`+d.ID()+`/`)

	return template.HTML(content)
}
