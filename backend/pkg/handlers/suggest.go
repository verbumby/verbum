package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Suggest handles _suggest http request
func Suggest(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	qbytes, err := json.Marshal(q)
	if err != nil {
		log.Println(errors.Wrap(err, "marshal q"))
		return
	}
	query := `{
		"_source": false,
		"suggest": {
			"HeadwordSuggest": {
				"prefix": ` + string(qbytes) + `,
				"completion": {
					"field": "Suggest",
					"skip_duplicates": true,
					"size": 10
				}
			}
		}
	}`

	url := viper.GetString("elastic.addr") + "/dict-*/_search"
	resp, err := http.Post(url, "application/json", strings.NewReader(query))
	if err != nil {
		log.Println(errors.Wrap(err, "query elastic"))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respbytes, _ := ioutil.ReadAll(resp.Body)
		log.Println(fmt.Errorf("query elastic: expected %d, got %d: %s", http.StatusOK, resp.StatusCode, string(respbytes)))
		return
	}

	respdata := struct {
		Suggest struct {
			HeadwordSuggest []struct {
				Options []struct {
					Text string `json:"text"`
				} `json:"options"`
			}
		} `json:"suggest"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&respdata); err != nil {
		log.Println(errors.Wrap(err, "unmarshal elastic resp"))
		return
	}

	data := []string{}
	for _, hws := range respdata.Suggest.HeadwordSuggest {
		for _, opt := range hws.Options {
			data = append(data, opt.Text)
		}
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println(errors.Wrap(err, "encode response"))
	}
}
