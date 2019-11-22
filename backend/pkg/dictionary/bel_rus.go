package dictionary

import "html/template"

import "strings"

// Rvblr dictionary
type BelRus struct {
	id    string
	title string
	slug  string
}

// ID implements Dictionary interface
func (d BelRus) ID() string {
	return d.id
}

// Title implements Dictionary interface
func (d BelRus) Title() string {
	return d.title
}

// Slug implements Dictionary interface
func (d BelRus) Slug() string {
	return d.slug
}

// ToHTML implements Dictionary interface
func (d BelRus) ToHTML(content string) template.HTML {
	content = strings.ReplaceAll(content, "<k>", "<p><v-hw>")
	content = strings.ReplaceAll(content, "</k>", "</v-hw></p>")
	return template.HTML(content)
}
