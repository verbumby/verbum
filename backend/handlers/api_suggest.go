package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/verbumby/verbum/backend/chttp"
	"github.com/verbumby/verbum/backend/dictionary"
	"github.com/verbumby/verbum/backend/htmlui"
	"github.com/verbumby/verbum/backend/storage"
	"github.com/verbumby/verbum/backend/textutil"
)

// suggestFuzzyMinLength is the elastic completion suggester min_length default:
// shorter inputs are never matched fuzzily, so there is no point in retrying them
const suggestFuzzyMinLength = 3

// suggestRequest builds the completion suggester query. Zero fuzziness means an
// exact prefix match.
func suggestRequest(q string, fuzziness int) map[string]any {
	completion := map[string]any{
		"field":           "Suggest",
		"skip_duplicates": true,
		"size":            10,
	}

	if fuzziness > 0 {
		completion["fuzzy"] = map[string]any{
			"fuzziness": fuzziness,
			// count the edit distance in code points: measured in bytes, a
			// single inserted, dropped or transposed cyrillic letter already
			// costs two edits and stays out of reach of fuzziness 1
			"unicode_aware": true,
		}
	}

	return map[string]any{
		"_source": false,
		"suggest": map[string]any{
			"HeadwordSuggest": map[string]any{
				"prefix":     q,
				"completion": completion,
			},
		},
	}
}

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

	// start with the exact prefix match and, if it yields nothing, retry with
	// increasingly tolerant fuzzy matching
	fuzziness := []int{0, 1, 2}
	if utf8.RuneCountInString(q) < suggestFuzzyMinLength {
		fuzziness = fuzziness[:1]
	}

	data := []string{}
	for _, f := range fuzziness {
		respbody := struct {
			Suggest struct {
				HeadwordSuggest []struct {
					Options []struct {
						Text string `json:"text"`
					} `json:"options"`
				}
			} `json:"suggest"`
		}{}

		if err := storage.Post("/"+inDictsStr+"/_search", suggestRequest(q, f), &respbody); err != nil {
			return fmt.Errorf("query elastic: %w", err)
		}

		for _, hws := range respbody.Suggest.HeadwordSuggest {
			for _, opt := range hws.Options {
				data = append(data, opt.Text)
			}
		}

		if len(data) > 0 {
			break
		}
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("encode response: %w", err)
	}

	return nil
}
