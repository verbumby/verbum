package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/verbumby/verbum/backend/pkg/app"
	"github.com/verbumby/verbum/backend/pkg/chttp"
	"github.com/verbumby/verbum/backend/pkg/db"
	"github.com/verbumby/verbum/backend/pkg/dict"
	"github.com/verbumby/verbum/backend/pkg/tm"

	"gopkg.in/reform.v1"
)

// RecordSaveHandler record create handler
type RecordSaveHandler struct {
	Table     reform.Table
	DB        *reform.DB
	AfterSave func(reform.Struct) error
}

func (h *RecordSaveHandler) ServeHTTP(w http.ResponseWriter, ctx *chttp.Context) error {
	record := h.Table.NewRecord()
	if err := json.NewDecoder(ctx.R.Body).Decode(record); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	var err error
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

	if h.AfterSave != nil {
		if err := h.AfterSave(record); err != nil {
			return errors.Wrap(err, "save record")
		}
	}
	return nil
}

// RecordListHandler record list handler
type RecordListHandler struct {
	Table   reform.Table
	DB      *reform.DB
	Filters []app.Filter
}

func (h *RecordListHandler) ServeHTTP(w http.ResponseWriter, ctx *chttp.Context) error {
	ffroms := []string{}
	fwheres := []string{}
	ffromArgss := []interface{}{}
	fwhereArgss := []interface{}{}
	for _, fmeta := range h.Filters {
		fgetname := "filter$" + fmeta.Name()
		if _, ok := ctx.R.URL.Query()[fgetname]; !ok {
			continue
		}
		ffrom, ffromArgs, fwhere, fwhereArgs, err := fmeta.ToSQL(ctx.R.URL.Query().Get(fgetname))
		if err != nil {
			http.Error(w, errors.Wrapf(err, "parse %s", fgetname).Error(), http.StatusBadRequest)
			return nil
		}
		if ffrom != "" {
			ffroms = append(ffroms, ffrom)
			ffromArgss = append(ffromArgss, ffromArgs...)
		}
		if fwhere != "" {
			fwheres = append(fwheres, fwhere)
			fwhereArgss = append(fwhereArgss, fwhereArgs...)
		}
	}

	queryTail := ""
	args := []interface{}{}
	if len(ffroms) > 0 {
		queryTail += strings.Join(ffroms, " ") + " "
		args = append(args, ffromArgss...)
	}
	if len(fwheres) > 0 {
		queryTail += "WHERE " + strings.Join(fwheres, " AND ") + " "
		args = append(args, fwhereArgss...)
	}

	limit := 20
	limitStr := ctx.R.URL.Query().Get("limit")
	if limitStr != "" {
		limit64, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil {
			http.Error(w, errors.Wrap(err, "parse limit query param").Error(), http.StatusBadRequest)
			return nil
		}
		limit = int(limit64)
	}
	args = append(args, limit)

	offset := 0
	offsetStr := ctx.R.URL.Query().Get("offset")
	if offsetStr != "" {
		offset64, err := strconv.ParseInt(offsetStr, 10, 32)
		if err != nil {
			http.Error(w, errors.Wrap(err, "parse offset query param").Error(), http.StatusBadRequest)
			return nil
		}
		offset = int(offset64)
	}
	args = append(args, offset)

	queryTail = queryTail + fmt.Sprintf(
		"LIMIT %s OFFSET %s",
		db.DB.Placeholder(1),
		db.DB.Placeholder(2),
	)
	records, err := h.DB.SelectAllFrom(h.Table, queryTail, args...)
	if err != nil {
		return errors.Wrap(err, "select from db")
	}

	json.NewEncoder(w).Encode(struct {
		Data interface{}
	}{
		Data: records,
	})
	return nil
}

// RecordFetchHandler record fetch handler
type RecordFetchHandler struct {
	ModelMeta app.ModelMeta
	DB        *reform.DB
}

func (h *RecordFetchHandler) ServeHTTP(w http.ResponseWriter, ctx *chttp.Context) error {
	vars := mux.Vars(ctx.R)
	id, err := parseInt(vars["ID"])
	if err != nil {
		http.Error(w, errors.Wrap(err, "parse ID param").Error(), http.StatusBadRequest)
		return nil
	}

	record := h.ModelMeta.NewModel()
	if err := h.DB.FindByPrimaryKeyTo(record, id); err != nil {
		return errors.Wrap(err, "find by primary key")
	}
	if err := record.LoadRelationships(); err != nil {
		return errors.Wrap(err, "load relationships")
	}

	json.NewEncoder(w).Encode(struct {
		Data interface{}
	}{
		Data: record,
	})
	return nil
}

func parseInt(str string) (int, error) {
	id64, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, errors.Wrap(err, "parse int")
	}

	return int(id64), nil
}

// IndexHandler serve admin index page
func IndexHandler(w http.ResponseWriter, ctx *chttp.Context) error {
	dicts, err := db.DB.SelectAllFrom(dict.DictTable, "")
	if err != nil {
		return errors.Wrap(err, "select all dicts")
	}

	data := struct {
		Dicts     []reform.Struct
		Principal *chttp.Principal
	}{
		Dicts:     dicts,
		Principal: ctx.P,
	}
	if err := tm.Render("admin", w, data); err != nil {
		return errors.Wrap(err, "render admin")
	}
	return nil
}
