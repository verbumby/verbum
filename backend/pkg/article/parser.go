package article

import (
	"fmt"
	"io"
	"strings"

	"github.com/verbumby/verbum/backend/pkg/fts"
	"github.com/verbumby/verbum/backend/pkg/typeahead"
	reform "gopkg.in/reform.v1"

	"github.com/verbumby/verbum/backend/pkg/headword"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

type parser struct {
	a   *Article
	hws [][]html.Token
}

func (p *parser) parse() error {
	a := p.a
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
			if err := p.validateTagTokens(queue[startingTagIDX:]); err != nil {
				return errors.Wrap(err, "processing tags sequence")
			}

			p.onTag(queue[startingTagIDX:])
		}
	}

	if err := p.indexHeadwords(); err != nil {
		return errors.Wrap(err, "update headwords")
	}

	return nil
}

func (p *parser) onTag(tokens []html.Token) {
	switch tokens[0].Data {
	case "v-hw":
		p.onHeadwordTag(tokens)
	}
}

func (p *parser) onHeadwordTag(tokens []html.Token) {
	p.hws = append(p.hws, tokens)
}

func (p *parser) validateTagTokens(tokens []html.Token) error {
	if len(tokens) < 2 {
		return fmt.Errorf("expected at least two (openning and closing) tokens")
	}

	qs := tokens[0]
	qe := tokens[len(tokens)-1]

	if qs.Type != html.StartTagToken {
		return fmt.Errorf("expected starting token type to be %s, got %s", html.StartTagToken, qs.Type)
	}

	if tokens[len(tokens)-1].Type != html.EndTagToken {
		return fmt.Errorf("expected ending token type to be %s, got %s", html.EndTagToken, qe.Type)
	}

	// check whether start and end tokens match
	if qs.Data != qe.Data {
		return fmt.Errorf("closing tag %s doesn't match opening tag %s", qe.Data, qs.Data)
	}

	return nil
}

func (p *parser) indexHeadwords() error {
	hws := []reform.Record{}
	ths := []reform.Record{}

	for i, tokens := range p.hws {
		if i >= 1<<4 {
			return fmt.Errorf("headword count exceeded allowed %d count", 1<<4)
		}
		idx := p.a.ID<<4 + int32(i)

		content := ""
		for _, t := range tokens[1 : len(tokens)-1] {
			content += t.Data
		}

		hws = append(hws, &headword.Headword{
			ID:        idx,
			Headword:  content,
			ArticleID: p.a.ID,
		})

		ths = append(ths, &typeahead.Typeahead{
			ID:        idx,
			Typeahead: content,
			ArticleID: p.a.ID,
		})
	}

	if len(hws) == 0 {
		return nil
	}

	for _, records := range [][]reform.Record{hws, ths} {
		tx, err := fts.Sphinx.Begin()
		if err != nil {
			return errors.Wrap(err, "begin sphinx tx")
		}

		table := records[0].Table()
		if err := replaceInto(records, tx); err != nil {
			tx.Rollback()
			return errors.Wrapf(err, "replace $s of article %d", table.Name(), p.a.ID)
		}

		firstRecord := records[0].PKValue()
		lastRecord := records[len(hws)-1].PKValue()
		if err := deleteWhereArticleID(table, p.a.ID, firstRecord, lastRecord, tx); err != nil {
			tx.Rollback()
			return errors.Wrapf(err, "delete obsolete headwords of article %d", p.a.ID)
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			return errors.Wrap(err, "commit sphinx tx")
		}
	}

	return nil
}

func replaceInto(records []reform.Record, tx reform.TXInterface) error {
	table := records[0].Table()
	columns := strings.Join(table.Columns(), ", ")

	placeholders := strings.Join(fts.Sphinx.Placeholders(1, len(table.Columns())), ", ")
	placeholderRows := []string{}
	for range records {
		placeholderRows = append(placeholderRows, "("+placeholders+")")
	}

	query := "REPLACE INTO " + table.Name() + " (" + columns + ") VALUES " +
		strings.Join(placeholderRows, ", ")

	values := []interface{}{}
	for _, record := range records {
		values = append(values, record.Values()...)
	}

	_, err := tx.Exec(query, values...)
	return err
}

func deleteWhereArticleID(table reform.Table, articleID, idsFirst, idLast interface{}, tx reform.TXInterface) error {
	query := "DELETE FROM `" + table.Name() + "` WHERE `article_id` = ? AND id NOT BETWEEN ? AND ?"
	_, err := tx.Exec(query, articleID, idsFirst, idLast)
	return err
}
