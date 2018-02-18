package article

import (
	"github.com/pkg/errors"
	"gopkg.in/reform.v1"
)

//go:generate reform

// Article represents an article
//
// reform:articles
type Article struct {
	ID      int32  `reform:"id,pk" json:"id"`
	Content string `reform:"content" json:"content"`
}

// Index updates sphinx index
func Index(record reform.Struct) error {
	p := parser{
		a: record.(*Article),
	}
	if err := p.parse(); err != nil {
		return errors.Wrap(err, "parse article")
	}
	return nil
}
