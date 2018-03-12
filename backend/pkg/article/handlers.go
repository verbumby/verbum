package article

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/verbumby/verbum/backend/pkg/chttp"
	"github.com/verbumby/verbum/backend/pkg/fts"
	"github.com/verbumby/verbum/backend/pkg/headword"
	"github.com/verbumby/verbum/backend/pkg/typeahead"
	reform "gopkg.in/reform.v1"
)

// RecordSaveHandler record create handler
type RecordSaveHandler struct {
	DB *reform.DB
}

func (h *RecordSaveHandler) ServeHTTP(w http.ResponseWriter, ctx *chttp.Context) error {
	record := ArticleTable.NewRecord()
	if err := json.NewDecoder(ctx.R.Body).Decode(record); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	article := record.(*Article)
	tokenGroups, err := parseArticle(article.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	if len(tokenGroups) == 0 {
		http.Error(w, fmt.Sprintf("no headwords, the article won't be searchable"), http.StatusBadRequest)
		return nil
	}

	article.Title = ""
	for _, t := range tokenGroups[0][1 : len(tokenGroups[0])-1] {
		article.Title += t.Data
	}

	if record.HasPK() {
		err = h.DB.Update(record)
		if err == reform.ErrNoRows {
			err = nil
		}
	} else {
		err = h.DB.Insert(record)
	}

	if err != nil {
		return errors.Wrap(err, "save record")
	}

	hws := []reform.Record{}
	ths := []reform.Record{}

	for i, tokens := range tokenGroups {
		if i >= 1<<4 {
			http.Error(w, fmt.Sprintf("headword count exceeded allowed %d count", 1<<4), http.StatusBadRequest)
			return nil
		}
		idx := article.ID<<4 + int32(i)

		content := ""
		for _, t := range tokens[1 : len(tokens)-1] {
			content += t.Data
		}

		hws = append(hws, &headword.Headword{
			ID:        idx,
			Headword:  content,
			ArticleID: article.ID,
		})

		ths = append(ths, &typeahead.Typeahead{
			ID:        idx,
			Typeahead: content,
			ArticleID: article.ID,
		})
	}

	if err := indexHeadwords(article.ID, hws); err != nil {
		return errors.Wrap(err, "index headwords")
	}
	if err := indexHeadwords(article.ID, ths); err != nil {
		return errors.Wrap(err, "index typeaheads")
	}

	return nil
}

func indexHeadwords(articleID int32, records []reform.Record) error {
	tx, err := fts.Sphinx.Begin()
	if err != nil {
		return errors.Wrap(err, "begin sphinx tx")
	}

	table := records[0].Table()
	if err := replaceInto(records, tx); err != nil {
		tx.Rollback()
		return errors.Wrapf(err, "replace $s of article %d", table.Name(), articleID)
	}

	firstRecord := records[0].PKValue()
	lastRecord := records[len(records)-1].PKValue()
	if err := deleteWhereArticleID(table, articleID, firstRecord, lastRecord, tx); err != nil {
		tx.Rollback()
		return errors.Wrapf(err, "delete obsolete headwords of article %d", articleID)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "commit sphinx tx")
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
