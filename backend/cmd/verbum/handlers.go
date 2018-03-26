package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
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
	Table reform.Table
	DB    *reform.DB
}

func (h *RecordListHandler) ServeHTTP(w http.ResponseWriter, ctx *chttp.Context) error {
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

	queryTail := fmt.Sprintf(
		"LIMIT %s OFFSET %s",
		db.DB.Placeholder(1),
		db.DB.Placeholder(2),
	)
	records, err := h.DB.SelectAllFrom(h.Table, queryTail, limit, offset)
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
	Table reform.Table
	DB    *reform.DB
}

func (h *RecordFetchHandler) ServeHTTP(w http.ResponseWriter, ctx *chttp.Context) error {
	vars := mux.Vars(ctx.R)
	id, err := parseInt(vars["ID"])
	if err != nil {
		http.Error(w, errors.Wrap(err, "parse ID param").Error(), http.StatusBadRequest)
		return nil
	}
	record, err := h.DB.FindByPrimaryKeyFrom(h.Table, id)
	if err != nil {
		return errors.Wrap(err, "find by primary key")
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
