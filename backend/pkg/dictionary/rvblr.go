package dictionary

import (
	"html/template"

	"gopkg.in/russross/blackfriday.v2"
)

// Rvblr dictionary
type Rvblr struct {
	id    string
	title string
	slug  string
}

// ID implements Dictionary interface
func (d Rvblr) ID() string {
	return d.id
}

// Title implements Dictionary interface
func (d Rvblr) Title() string {
	return d.title
}

// Slug implements Dictionary interface
func (d Rvblr) Slug() string {
	return d.slug
}

// ToHTML implements Dictionary interface
func (d Rvblr) ToHTML(content string) template.HTML {
	return template.HTML(blackfriday.Run([]byte(content)))
}
