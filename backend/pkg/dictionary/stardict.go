package dictionary

import (
	"html/template"
	"strings"
)

// Stardict dictionary
type Stardict struct {
	id      string
	indexID string
	aliases []string
	title   string
}

// ID implements Dictionary interface
func (d Stardict) ID() string {
	return d.id
}

func (d Stardict) IndexID() string {
	if d.indexID == "" {
		return d.id
	}
	return d.indexID
}

func (d Stardict) Aliases() []string {
	return d.aliases
}

// Title implements Dictionary interface
func (d Stardict) Title() string {
	return d.title
}

// ToHTML implements Dictionary interface
func (d Stardict) ToHTML(content, title string) template.HTML {
	content = strings.ReplaceAll(content, "<k>", "<p><v-hw>")
	content = strings.ReplaceAll(content, "</k>", "</v-hw></p>")
	return template.HTML(content)
}
