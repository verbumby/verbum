package dictionary

import (
	"html/template"

	"gopkg.in/russross/blackfriday.v2"
)

// Markdown dictionary
type Markdown struct {
	id    string
	title string
	slug  string
}

// ID implements Dictionary interface
func (d Markdown) ID() string {
	return d.id
}

// Title implements Dictionary interface
func (d Markdown) Title() string {
	return d.title
}

// Slug implements Dictionary interface
func (d Markdown) Slug() string {
	return d.slug
}

// ToHTML implements Dictionary interface
func (d Markdown) ToHTML(content, title string) template.HTML {
	return template.HTML(blackfriday.Run([]byte(content)))
}
