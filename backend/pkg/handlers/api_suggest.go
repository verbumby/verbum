package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/verbumby/verbum/backend/pkg/chttp"
	"github.com/verbumby/verbum/backend/pkg/storage"
)

// APISuggest handles suggest http request
func APISuggest(w http.ResponseWriter, rctx *chttp.Context) error {
	q := rctx.R.URL.Query().Get("q")
	if len(q) > 500 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil
	}

	reqbody := map[string]interface{}{
		"_source": false,
		"suggest": map[string]interface{}{
			"HeadwordSuggest": map[string]interface{}{
				"prefix": q,
				"completion": map[string]interface{}{
					"field":           "Suggest",
					"skip_duplicates": true,
					"size":            10,
				},
			},
		},
	}

	respbody := struct {
		Suggest struct {
			HeadwordSuggest []struct {
				Options []struct {
					Text string `json:"text"`
				} `json:"options"`
			}
		} `json:"suggest"`
	}{}

	if err := storage.Post("/dict-*/_search", reqbody, &respbody); err != nil {
		return fmt.Errorf("query elastic: %w", err)
	}

	data := []string{}
	for _, hws := range respbody.Suggest.HeadwordSuggest {
		for _, opt := range hws.Options {
			data = append(data, opt.Text)
		}
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("encode response: %w", err)
	}

	return nil
}
