package dictionary

import (
	"html/template"
	"regexp"

	"golang.org/x/text/unicode/norm"
)

type HTML struct {
	id      string
	indexID string
	aliases []string
	title   string
}

func (d HTML) ID() string {
	return d.id
}

func (d HTML) IndexID() string {
	if d.indexID == "" {
		return d.id
	}
	return d.indexID
}

func (d HTML) Aliases() []string {
	return d.aliases
}

func (d HTML) Title() string {
	return d.title
}

func (d HTML) ToHTML(content, title string) template.HTML {
	content = norm.NFD.String(content)

	var re = regexp.MustCompile(`.\x{0301}`)
	var substitution = "<v-accent>$0</v-accent>"

	content = re.ReplaceAllString(content, substitution)

	content = norm.NFC.String(content)

	return template.HTML(content)
}
