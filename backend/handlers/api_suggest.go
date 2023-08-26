package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/verbumby/verbum/backend/chttp"
	"github.com/verbumby/verbum/backend/dictionary"
	"github.com/verbumby/verbum/backend/htmlui"
	"github.com/verbumby/verbum/backend/storage"
	"github.com/verbumby/verbum/backend/textutil"
)

// APISuggest handles suggest http request
func APISuggest(w http.ResponseWriter, rctx *chttp.Context) error {
	urlQuery := htmlui.Query([]htmlui.QueryParam{
		htmlui.NewStringQueryParam("q", ""),
		htmlui.NewInDictsQueryParam("in"),
	})
	urlQuery.From(rctx.R.URL.Query())

	q := urlQuery.Get("q").(*htmlui.StringQueryParam).Value()
	inDicts := urlQuery.Get("in").(*htmlui.InDictsQueryParam).Value()
	if len(q) > 1000 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil
	}
	q = textutil.NormalizeQuery(q)

	inDictsStr := ""
	for _, d := range inDicts {
		if len(inDictsStr) == 0 {
			inDictsStr = "sugg-" + dictionary.Get(d).IndexID()
		} else {
			inDictsStr += ",sugg-" + dictionary.Get(d).IndexID()
		}
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

	if err := storage.Post("/"+inDictsStr+"/_search", reqbody, &respbody); err != nil {
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
