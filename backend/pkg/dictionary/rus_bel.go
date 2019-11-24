package dictionary

import "html/template"

import "strings"

// RusBel dictionary
type RusBel struct {
	id    string
	title string
	slug  string
}

// ID implements Dictionary interface
func (d RusBel) ID() string {
	return d.id
}

// Title implements Dictionary interface
func (d RusBel) Title() string {
	return d.title
}

// Slug implements Dictionary interface
func (d RusBel) Slug() string {
	return d.slug
}

// ToHTML implements Dictionary interface
func (d RusBel) ToHTML(content string) template.HTML {
	content = strings.ReplaceAll(content, "<k>", "<p><v-hw>")
	content = strings.ReplaceAll(content, "</k>", "</v-hw></p>")
	return template.HTML(content)
}
