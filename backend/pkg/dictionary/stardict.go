package dictionary

import "html/template"

import "strings"

// Stardict dictionary
type Stardict struct {
	id    string
	title string
	slug  string
}

// ID implements Dictionary interface
func (d Stardict) ID() string {
	return d.id
}

// Title implements Dictionary interface
func (d Stardict) Title() string {
	return d.title
}

// Slug implements Dictionary interface
func (d Stardict) Slug() string {
	return d.slug
}

// ToHTML implements Dictionary interface
func (d Stardict) ToHTML(content, title string) template.HTML {
	content = strings.ReplaceAll(content, "<k>", "<p><v-hw>")
	content = strings.ReplaceAll(content, "</k>", "</v-hw></p>")
	return template.HTML(content)
}
