package dictionary

import (
	"html/template"

	"github.com/verbumby/verbum/backend/pkg/dictionary/dslparser"
)

// DSL dictionary
type DSL struct {
	id                    string
	title                 string
	slug                  string
	includeTitleInContent bool
}

// ID implements Dictionary interface
func (d DSL) ID() string {
	return d.id
}

// Title implements Dictionary interface
func (d DSL) Title() string {
	return d.title
}

// Slug implements Dictionary interface
func (d DSL) Slug() string {
	return d.slug
}

// ToHTML implements Dictionary interface
func (d DSL) ToHTML(content, title string) template.HTML {
	htmlv, err := dslparser.DSLParser(content)
	if err != nil {
		panic(err)
	}
	if d.includeTitleInContent {
		htmlv = `<p><v-hw>` + title + `</v-hw></p>` + htmlv
	}
	return template.HTML(htmlv)
}
