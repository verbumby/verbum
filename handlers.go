package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"gopkg.in/reform.v1"
)

// RecordSaveHandler record create handler
type RecordSaveHandler struct {
	Table reform.Table
	DB    *reform.DB
}

func (h *RecordSaveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	record := h.Table.NewRecord()
	if err := json.NewDecoder(r.Body).Decode(record); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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
		log.Printf("save record: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// RecordListHandler record list handler
type RecordListHandler struct {
	Table reform.Table
	DB    *reform.DB
}

func (h *RecordListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := "SELECT %s FROM %s"
	query = fmt.Sprintf(query, strings.Join(h.Table.Columns(), ", "), h.Table.Name())
	records, err := h.DB.SelectAllFrom(h.Table, "")
	if err != nil {
		log.Printf("list record: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(struct {
		Data interface{} `json:"data"`
	}{
		Data: records,
	})
}

// RecordFetchHandler record fetch handler
type RecordFetchHandler struct {
	Table reform.Table
	DB    *reform.DB
}

func (h *RecordFetchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := parseIntID(vars["ID"])
	if err != nil {
		http.Error(w, errors.Wrap(err, "parse ID param").Error(), http.StatusBadRequest)
		return
	}
	record, err := h.DB.FindByPrimaryKeyFrom(h.Table, id)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		log.Println(errors.Wrap(err, "find by primary key"))
		return
	}

	json.NewEncoder(w).Encode(struct {
		Data interface{} `json:"data"`
	}{
		Data: record,
	})
}

func parseIntID(str string) (int, error) {
	id64, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, errors.Wrap(err, "parse int")
	}

	return int(id64), nil
}
