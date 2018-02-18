package article

import (
	"github.com/pkg/errors"
)

//go:generate reform

// Article represents an article
//
// reform:articles
type Article struct {
	ID      int32  `reform:"id,pk" json:"id"`
	Content string `reform:"content" json:"content"`
}

// BeforeInsert implements reform.BeforeInserter
func (a *Article) BeforeInsert() error {
	return a.BeforeSave()
}

// BeforeUpdate implements reform.BeforeUpdater
func (a *Article) BeforeUpdate() error {
	return a.BeforeSave()
}

// BeforeSave updates sphinx index/indices
func (a *Article) BeforeSave() error {
	p := parser{
		a: a,
	}
	if err := p.parse(); err != nil {
		return errors.Wrap(err, "parse article")
	}
	return nil
}
