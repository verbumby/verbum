package dictionary

import (
	"fmt"
	"html/template"

	"github.com/verbumby/verbum/backend/dictionary/dslparser"
)

// DSL dictionary
type DSL struct {
	Common
	RenderRefs bool
}

// ToHTML implements Dictionary interface
func (d DSL) ToHTML(content string) template.HTML {
	htmlvitf, err := dslparser.Parse(
		"article",
		[]byte(content),
		dslparser.GlobalStore("dictID", d.ID()),
		dslparser.GlobalStore("renderRefs", d.RenderRefs),
	)
	if err != nil {
		panic(fmt.Errorf("parse article: %w", err))
	}
	htmlv := htmlvitf.(string)
	if d.abbrevs != nil {
		htmlv = renderAbbrevs(htmlv, d.abbrevs)
	}

	htmlv = wrapAccentedChars(htmlv)

	return template.HTML(htmlv)
}
