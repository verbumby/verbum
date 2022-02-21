package dictionary

import (
	"html/template"
	"strings"
)

// Stardict dictionary
type Stardict struct {
	id    string
	title string
}

// ID implements Dictionary interface
func (d Stardict) ID() string {
	return d.id
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
