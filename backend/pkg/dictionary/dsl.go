package dictionary

import (
	"fmt"
	"html/template"

	"github.com/verbumby/verbum/backend/pkg/dictionary/dslparser"
)

// DSL dictionary
type DSL struct {
	id      string
	indexID string
	aliases []string
	title   string
	abbrevs map[string]string
}

// ID implements Dictionary interface
func (d DSL) ID() string {
	return d.id
}

func (d DSL) IndexID() string {
	if d.indexID == "" {
		return d.id
	}
	return d.indexID
}

func (d DSL) Aliases() []string {
	return d.aliases
}

// Title implements Dictionary interface
func (d DSL) Title() string {
	return d.title
}

// ToHTML implements Dictionary interface
func (d DSL) ToHTML(content string) template.HTML {
	htmlvitf, err := dslparser.Parse(
		"article",
		[]byte(content),
		dslparser.GlobalStore("dictID", d.ID()),
	)
	if err != nil {
		panic(fmt.Errorf("parse article: %w", err))
	}
	htmlv := htmlvitf.(string)
	if d.abbrevs != nil {
		htmlv = renderAbbrevs(htmlv, d.abbrevs)
	}
	return template.HTML(htmlv)
}
