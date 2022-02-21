package dictionary

import (
	"fmt"
	"html/template"

	"github.com/verbumby/verbum/backend/pkg/dictionary/dslparser"
)

// DSL dictionary
type DSL struct {
	id                    string
	title                 string
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

// ToHTML implements Dictionary interface
func (d DSL) ToHTML(content, title string) template.HTML {
	htmlvitf, err := dslparser.Parse(
		"article",
		[]byte(content),
		dslparser.GlobalStore("dictID", d.ID()),
	)
	if err != nil {
		panic(fmt.Errorf("parse article: %w", err))
	}
	htmlv := htmlvitf.(string)
	if d.includeTitleInContent {
		htmlv = `<p><v-hw>` + title + `</v-hw></p>` + htmlv
	}
	return template.HTML(htmlv)
}
