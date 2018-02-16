package article

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/verbumby/verbum/fts"

	"github.com/verbumby/verbum/headword"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
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

	queue := []html.Token{}
	z := html.NewTokenizer(strings.NewReader(a.Content))
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			err := z.Err()
			if err != io.EOF {
				return errors.Wrap(err, "tokenize html")
			}
			break
		}

		t := z.Token()
		queue = append(queue, t)
		if tt == html.EndTagToken {
			startingTagIDX := len(queue) - 1
			for startingTagIDX >= 0 && queue[startingTagIDX].Type != html.StartTagToken {
				startingTagIDX--
			}
			if startingTagIDX < 0 {
				return fmt.Errorf("no starting tag found for closing tag %s", t.Type)
			}
			if err := processTokenQueue(queue[startingTagIDX:]); err != nil {
				return errors.Wrap(err, "processing tags sequence")
			}
		}

	}
	content := ""
	for _, t := range queue {
		content += t.String()
	}
	a.Content = content
	return nil
}

func processTokenQueue(q []html.Token) error {
	if len(q) < 2 {
		return fmt.Errorf("expected at least two (openning and closing) tokens")
	}

	qs := q[0]
	qe := q[len(q)-1]

	if qs.Type != html.StartTagToken {
		return fmt.Errorf("expected starting token type to be %s, got %s", html.StartTagToken, qs.Type)
	}

	if q[len(q)-1].Type != html.EndTagToken {
		return fmt.Errorf("expected ending token type to be %s, got %s", html.EndTagToken, qe.Type)
	}

	// check whether start and end tokens match
	if qs.Data != qe.Data {
		return fmt.Errorf("closing tag %s doesn't match opening tag %s", qe.Data, qs.Data)
	}

	if qs.Data != "v-hw" {
		return nil
	}

	hwcontent := ""
	for _, t := range q[1 : len(q)-1] {
		hwcontent += t.String()
	}

	hw := &headword.Headword{
		Headword: hwcontent,
	}

	var idxAttr *html.Attribute
	for i := range qs.Attr {
		if qs.Attr[i].Key == "idx" {
			idxAttr = &qs.Attr[i]
		}
	}
	if idxAttr == nil {
		qs.Attr = append(qs.Attr, html.Attribute{Key: "idx"})
		idxAttr = &qs.Attr[len(qs.Attr)-1]
	}

	tx, err := fts.Sphinx.Begin()
	if err != nil {
		return errors.Wrap(err, "begin sphinx tx")
	}

	if idxAttr.Val != "" {
		idx64, err := strconv.ParseInt(idxAttr.Val, 10, 32)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("%s is not valid id of a headword", idxAttr.Val)
		}

		hw.ID = int32(idx64)
	} else {
		if err := tx.QueryRow("SELECT MAX(id) FROM headwords").Scan(&hw.ID); err != nil {
			tx.Rollback()
			return errors.Wrap(err, "select max headword id")
		}

		hw.ID++
		idxAttr.Val = strconv.FormatInt(int64(hw.ID), 10)
	}

	columns := hw.Table().Columns()
	placeholders := fts.Sphinx.Placeholders(1, len(columns))
	query := "REPLACE INTO headwords (" + strings.Join(columns, ", ") + ") " +
		"VALUES (" + strings.Join(placeholders, ", ") + ")"
	if _, err := fts.Sphinx.Exec(query, hw.Values()...); err != nil {
		tx.Rollback()
		return errors.Wrapf(err, "replace headword %d", hw.ID)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return errors.Wrapf(err, "commit sphinx tx after headowrd %d replace", hw.ID)
	}

	return nil
}
