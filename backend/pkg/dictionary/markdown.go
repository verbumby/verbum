package dictionary

import (
	"html/template"

	"gopkg.in/russross/blackfriday.v2"
)

// Markdown dictionary
type Markdown struct {
	id      string
	indexID string
	aliases []string
	title   string
}

// ID implements Dictionary interface
func (d Markdown) ID() string {
	return d.id
}

func (d Markdown) IndexID() string {
	if d.indexID == "" {
		return d.id
	}
	return d.indexID
}

func (d Markdown) Aliases() []string {
	return d.aliases
}

// Title implements Dictionary interface
func (d Markdown) Title() string {
	return d.title
}

// ToHTML implements Dictionary interface
func (d Markdown) ToHTML(content string) template.HTML {
	return template.HTML(blackfriday.Run([]byte(content)))
}
