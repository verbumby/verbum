package dictionary

import (
	"html/template"

	"github.com/verbumby/verbum/backend/pkg/dictionary/krapivaparser"
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
	htmlv, err := krapivaparser.KrapivaParser(content)
	if err != nil {
		panic(err)
	}
	return template.HTML(htmlv)
}
