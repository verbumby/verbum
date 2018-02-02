package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"gopkg.in/reform.v1"
)

// RecordCreateHandler record create handler
type RecordCreateHandler struct {
	Table reform.Table
	DB    *reform.DB
}

func (h *RecordCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	record := h.Table.NewRecord()
	if err := json.NewDecoder(r.Body).Decode(record); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.DB.Save(record); err != nil {
		log.Printf("create dict: %v", err)
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
		log.Printf("list dict: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(struct {
		Data interface{} `json:"data"`
	}{
		Data: records,
	})
}
