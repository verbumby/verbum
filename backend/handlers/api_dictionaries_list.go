package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/verbumby/verbum/backend/chttp"
	"github.com/verbumby/verbum/backend/dictionary"
)

// APIDictionariesList dictionaries list endpoint
func APIDictionariesList(w http.ResponseWriter, rctx *chttp.Context) error {
	type dictview struct {
		ID         string
		Aliases    []string
		Title      string
		HasPreface bool
		HasAbbrevs bool
		ScanURL    string
	}

	resp := struct {
		Dicts    []dictview
		Sections []dictionary.Section
	}{
		Sections: dictionary.GetAllSections(),
	}

	for _, d := range dictionary.GetAll() {
		resp.Dicts = append(resp.Dicts, dictview{
			ID:         d.ID(),
			Aliases:    d.Aliases(),
			Title:      d.Title(),
			HasPreface: d.Preface() != "",
			HasAbbrevs: d.Abbrevs() != nil,
			ScanURL:    d.ScanURL(),
		})
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return fmt.Errorf("encode response: %w", err)
	}
	return nil
}
