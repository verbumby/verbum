package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/verbumby/verbum/backend/pkg/chttp"
	"github.com/verbumby/verbum/backend/pkg/storage"
)

// Suggest handles _suggest http request
func Suggest(w http.ResponseWriter, rctx *chttp.Context) error {
	q := rctx.R.URL.Query().Get("q")

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
		return errors.Wrap(err, "query elastic")
	}

	data := []string{}
	for _, hws := range respbody.Suggest.HeadwordSuggest {
		for _, opt := range hws.Options {
			data = append(data, opt.Text)
		}
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return errors.Wrap(err, "encode response")
	}

	return nil
}
