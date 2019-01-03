package esmigrate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/verbumby/verbum/backend/pkg/article"
	"github.com/verbumby/verbum/backend/pkg/db"
	"github.com/verbumby/verbum/backend/pkg/dict"
)

// GetCmd returns cobra cmd
func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use: "es-migrate",
		Run: Run,
	}
}

// Run executes the command
func Run(cmd *cobra.Command, args []string) {
	err := func() error {
		if err := createRvblrIndex(); err != nil {
			return errors.Wrap(err, "create rvblr index")
		}
		if err := indexRvblr(); err != nil {
			return errors.Wrap(err, "index rvblr")
		}
		return nil
	}()
	if err != nil {
		log.Fatal(err)
	}

	structs, err := db.DB.SelectAllFrom(dict.DictTable, "")
	if err != nil {
		log.Fatalf("select all dicts: %v", err)
	}
	for _, s := range structs {
		d := s.(*dict.Dict)
		fmt.Println(d)
	}
}

func createRvblrIndex() error {
	q := `{
		"settings": {
			"number_of_shards" : 1,
			"number_of_replicas" : 0,
			"analysis": {
				"analyzer": {
					"hw": {
						"type": "custom",
						"tokenizer": "keyword",
						"filter": ["lowercase"]
					},
					"hw_smaller": {
						"type": "custom",
						"tokenizer": "hw_smaller",
						"filter": ["lowercase"]
					}
				},
				"tokenizer": {
					"hw_smaller": {
						"type": "char_group",
						"tokenize_on_chars": ["-"]
					}
				}
			}
		},
		"mappings": {
			"_doc": {
				"properties": {
					"Headword": {
						"type": "text",
						"analyzer": "hw",
						"fields": {
							"Smaller": {
								"type": "text",
								"analyzer": "hw_smaller",
								"search_analyzer": "hw"
							}
						}
					},
					"HeadwordAlt": {
						"type": "text",
						"analyzer": "hw",
						"fields": {
							"Smaller": {
								"type": "text",
								"analyzer": "hw_smaller",
								"search_analyzer": "hw"
							}
						}
					},
					"Suggest": {
						"type": "completion",
						"analyzer": "hw"
					},
					"Content": {
						"type": "text",
						"index": false
					}
				}
			}
		}
	}`

	url := viper.GetString("elastic.addr") + "/dict-rvblr"
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(q))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return errors.Wrap(err, "new request")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	respBodyBytes, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected http status %d, got %d: %s", http.StatusOK, resp.StatusCode, string(respBodyBytes))
	}

	log.Printf("dict-rvblr index has been created: %s", string(respBodyBytes))

	return nil
}

func indexRvblr() error {
	offset := 0
	c := 0
	for {
		tail := fmt.Sprintf("LIMIT 100 OFFSET %d", offset)
		records, err := db.DB.SelectAllFrom(article.ArticleTable, tail)
		if err != nil {
			return errors.Wrap(err, "select articles")
		}

		if len(records) == 0 {
			break
		}
		buff := &bytes.Buffer{}
		for _, record := range records {
			a := record.(*article.Article)
			buff.WriteString(`{"index": {} }` + "\n")

			type suggestinputt struct {
				Input  string `json:"input"`
				Weight int    `json:"weight"`
			}
			type doct struct {
				Headword    []string
				HeadwordAlt []string
				Suggest     []suggestinputt
				Content     string
			}
			doc := doct{
				Headword:    []string{},
				HeadwordAlt: []string{},
				Suggest:     []suggestinputt{},
				Content:     a.Content,
			}

			hws, hwalts, err := article.RvblrParse(a.Content)
			if err != nil {
				return errors.Wrapf(err, "parse rvblr content: %s", "<content>")
			}

			for _, hw := range hws {
				doc.Headword = append(doc.Headword, hw)
				doc.Suggest = append(doc.Suggest, suggestinputt{
					Input:  hw,
					Weight: 4,
				})
			}

			for _, hw := range hwalts {
				doc.HeadwordAlt = append(doc.HeadwordAlt, hw)
				doc.Suggest = append(doc.Suggest, suggestinputt{
					Input:  hw,
					Weight: 2,
				})
			}

			if err := json.NewEncoder(buff).Encode(doc); err != nil {
				return errors.Wrapf(err, "json encode of doc %v", doc)
			}
			buff.WriteString("\n")
		}

		err = func() error {
			url := viper.GetString("elastic.addr") + "/dict-rvblr/_doc/_bulk"
			req, err := http.NewRequest(http.MethodPost, url, buff)
			req.Header.Add("Content-Type", "application/x-ndjson")
			if err != nil {
				return errors.Wrap(err, "new request")
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return errors.Wrap(err, "do request")
			}
			defer resp.Body.Close()

			respBodyBytes, _ := ioutil.ReadAll(resp.Body)
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("expected http status %d, got %d: %s", http.StatusOK, resp.StatusCode, string(respBodyBytes))
			}

			log.Printf("offset indexed: %d", offset)
			return nil
		}()

		if err != nil {
			return err
		}

		offset += 100
	}
	fmt.Println(c)
	return nil
}
