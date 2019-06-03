package dictionary

import (
	"html/template"
)

// Krapiva dictionary
type Krapiva struct {
	id    string
	title string
	slug  string
}

// ID implements Dictionary interface
func (d Krapiva) ID() string {
	return d.id
}

// Title implements Dictionary interface
func (d Krapiva) Title() string {
	return d.title
}

// Slug implements Dictionary interface
func (d Krapiva) Slug() string {
	return d.slug
}

// ToHTML implements Dictionary interface
func (d Krapiva) ToHTML(content string) template.HTML {
	return template.HTML(content)
}
