package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/verbumby/verbum/backend/pkg/chttp"
	"github.com/verbumby/verbum/backend/pkg/dictionary"
)

// APIDictionariesList dictionaries list endpoint
func APIDictionariesList(w http.ResponseWriter, rctx *chttp.Context) error {
	type dictview struct {
		ID      string
		Aliases []string
		Title   string
	}
	toencode := []dictview{}
	for _, d := range dictionary.GetAll() {
		toencode = append(toencode, dictview{ID: d.ID(), Aliases: d.Aliases(), Title: d.Title()})
	}
	if err := json.NewEncoder(w).Encode(toencode); err != nil {
		return fmt.Errorf("encode response: %w", err)
	}
	return nil
}
