package krapiva

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/verbumby/verbum/backend/pkg/textutil"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/verbumby/verbum/backend/pkg/storage"
)

type ParsedKrapiva struct {
	Preambles []interface{}
	Articles  []ParsedArticle
}

type ParsedPreamble struct {
	Key   string
	Value string
}

type ParsedArticle struct {
	Headwords []string
	Body      string
}

// Command creates cobra command
func Command() *cobra.Command {
	var filename string
	result := &cobra.Command{
		Use:   "krapiva-import",
		Short: "Imports krapiva dictionary",
		Long:  "Imports krapiva dictionary",
		Run: func(cmd *cobra.Command, args []string) {
			itf, err := ParseFile(filename)
			if err != nil {
				log.Fatalf("parser failed: %v", err)
			}
			d := itf.(ParsedKrapiva)
			log.Println("parsed dictionary file")
			// fmt.Printf("%v", krapiva)

			if err := createIndex(); err != nil {
				log.Fatalf("create index: %v", err)
			}
			log.Println("created index")

			if err := indexArticles(d.Articles); err != nil {
				log.Fatalf("index articles: %v", err)
			}
			// create index mapping and settings
			// process all articles
		},
	}

	result.PersistentFlags().StringVar(&filename, "filename", "", "filename of the dict")
	return result
}

func createIndex() error {
	respbody := map[string]interface{}{}
	err := storage.Put("/dict-krapiva", map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":  1,
			"max_result_window": 100000,
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"hw": map[string]interface{}{
						"filter":    []string{"lowercase"},
						"type":      "custom",
						"tokenizer": "keyword",
					},
					"hw_smaller": map[string]interface{}{
						"filter":    []string{"lowercase"},
						"type":      "custom",
						"tokenizer": "hw_smaller",
					},
				},
				"tokenizer": map[string]interface{}{
					"hw_smaller": map[string]interface{}{
						"type":              "char_group",
						"tokenize_on_chars": []string{"-", " ", "(", ")", ","},
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"_doc": map[string]interface{}{
				"properties": map[string]interface{}{
					"Title": map[string]interface{}{"type": "keyword"},
					"Headword": map[string]interface{}{
						"type":     "text",
						"analyzer": "hw",
						"fields": map[string]interface{}{
							"Smaller": map[string]interface{}{
								"type":            "text",
								"analyzer":        "hw_smaller",
								"search_analyzer": "hw",
							},
						},
					},
					"HeadwordAlt": map[string]interface{}{
						"type":     "text",
						"analyzer": "hw",
						"fields": map[string]interface{}{
							"Smaller": map[string]interface{}{
								"type":            "text",
								"analyzer":        "hw_smaller",
								"search_analyzer": "hw",
							},
						},
					},
					"Prefix": map[string]interface{}{
						"type": "nested",
						"properties": map[string]interface{}{
							"Letter1": map[string]interface{}{"type": "keyword"},
							"Letter2": map[string]interface{}{"type": "keyword"},
							"Letter3": map[string]interface{}{"type": "keyword"},
						},
					},
					"Suggest": map[string]interface{}{
						"type":                         "completion",
						"analyzer":                     "hw",
						"preserve_separators":          true,
						"preserve_position_increments": true,
						"max_input_length":             50,
					},
					"Content": map[string]interface{}{
						"type":  "text",
						"index": false,
					},
				},
			},
		},
	}, &respbody)
	if err != nil {
		return errors.Wrap(err, "storage put")
	}
	return nil
}

func indexArticles(as []ParsedArticle) error {
	idcache := map[string]int{}

	for _, a := range as {
		err := indexArticle(a, idcache)
		if err != nil {
			errors.Wrapf(err, "index article %v", a)
		}
		fmt.Print(".")
	}
	fmt.Println(" OK")
	return nil
}

func indexArticle(a ParsedArticle, idcache map[string]int) error {
	hws := []string{}
	hwsalt := []string{}
	suggests := []map[string]interface{}{}
	prefixes := []map[string]string{}

	var re = regexp.MustCompile(`(?m)^\t`)
	a.Body = re.ReplaceAllLiteralString(a.Body, "")
	bodylower := strings.ToLower(a.Body)

	for _, phw := range a.Headwords {
		phw = strings.TrimSpace(phw)
		phw = strings.ReplaceAll(phw, "\\(", "(")
		phw = strings.ReplaceAll(phw, "\\)", ")")

		phrase := fmt.Sprintf("[b][ex][lang id=1049][c steelblue]%s[/c][/lang][/ex][/b]", phw)
		phrase = strings.ToLower(phrase)

		phw = strings.ReplaceAll(phw, "...", "")

		if strings.Contains(bodylower, phrase) || strings.Contains(phw, "(") {
			hwsalt = append(hwsalt, phw)
			suggests = append(suggests, map[string]interface{}{
				"input":  phw,
				"weight": 2,
			})
		} else {
			hws = append(hws, phw)
			suggests = append(suggests, map[string]interface{}{
				"input":  phw,
				"weight": 4,
			})
			prefix := map[string]string{}
			i := 0
			for _, r := range phw {
				if i > 2 {
					break
				}
				prefix[fmt.Sprintf("Letter%d", i+1)] = string(r)
				i++
			}
			prefixes = append(prefixes, prefix)
		}
	}

	id := assembleID(hws)
	idcache[id]++
	if idcache[id] > 1 {
		id = fmt.Sprintf("%s-%d", id, idcache[id])
	}

	if err := storage.Put("/dict-krapiva/_doc/"+id, map[string]interface{}{
		"Title":       strings.Join(hws, ","),
		"Headword":    hws,
		"HeadwordAlt": hwsalt,
		"Suggest":     suggests,
		"Prefix":      prefixes,
		"Content":     a.Body,
	}, nil); err != nil {
		return errors.Wrap(err, "storage put")
	}

	//fmt.Println(id, "\t\t", strings.Join(hws, "; "), "\t\t", strings.Join(hwsalt, "; "))

	return nil
}

func assembleID(hws []string) string {
	romanized := []string{}
	for _, hw := range hws {
		romanized = append(romanized, textutil.RomanizeBelarusian(hw))
	}
	result := strings.Join(romanized, "-")
	return textutil.Slugify(result)
}
