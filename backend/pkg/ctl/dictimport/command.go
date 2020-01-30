package dictimport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/verbumby/verbum/backend/pkg/ctl/dictimport/dictparser"
	"github.com/verbumby/verbum/backend/pkg/ctl/dictimport/dictparser/dsl"
	"github.com/verbumby/verbum/backend/pkg/storage"
	"github.com/verbumby/verbum/backend/pkg/textutil"
)

// Command creates a cobra command
func Command() *cobra.Command {
	c := &commandController{}
	result := &cobra.Command{
		Use:   "import",
		Short: "Imports a dictionary",
		Long:  "Imports a dictionary",
		Run:   c.Run,
	}

	result.PersistentFlags().StringVar(&c.filename, "filename", "", "filename of the dict")
	result.PersistentFlags().StringVar(&c.format, "format", "", "dsl|stardict")
	result.PersistentFlags().StringVar(&c.indexID, "index-id", "", "storage index id")
	return result
}

type commandController struct {
	filename string
	format   string
	indexID  string
}

func (c *commandController) Run(cmd *cobra.Command, args []string) {
	fmt.Println(c.filename)
	if err := c.run(); err != nil {
		log.Fatal(err)
	}
}

func (c *commandController) run() error {
	f, err := os.Open(c.filename)
	if err != nil {
		return errors.Wrapf(err, "open %s file", c.filename)
	}
	defer f.Close()

	var d dictparser.Dictionary
	switch c.format {
	case "dsl":
		d, err = dsl.ParseDSLReader(c.filename, f)
	default:
		err = fmt.Errorf("unsupported format %s", c.format)
	}
	if err != nil {
		return errors.Wrap(err, "parse dictionary")
	}

	log.Printf("found %d articles in the dictionary", len(d.Articles))

	if err := c.createIndex(len(d.Articles)); err != nil {
		return errors.Wrap(err, "create index")
	}

	if err := c.indexArticles(d); err != nil {
		return errors.Wrap(err, "index articles")
	}

	return nil
}

func (c *commandController) createIndex(maxResultWindow int) error {
	err := storage.Put("/dict-"+c.indexID, map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":  1,
			"max_result_window": maxResultWindow,
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
	}, nil)
	if err != nil {
		return errors.Wrap(err, "storage put")
	}
	return nil
}

func (c *commandController) indexArticles(d dictparser.Dictionary) error {
	idcache := map[string]int{}

	buff := &bytes.Buffer{}
	for i, a := range d.Articles {
		suggests := []map[string]interface{}{}
		prefixes := []map[string]string{}

		for _, phw := range a.HeadwordsAlt {
			suggests = append(suggests, map[string]interface{}{
				"input":  phw,
				"weight": 2,
			})
		}

		for _, phw := range a.Headwords {
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

		// TODO: flag to control id assemble strategy
		id := c.assembleID(a.Headwords)
		idcache[id]++
		if idcache[id] > 1 {
			id = fmt.Sprintf("%s-%d", id, idcache[id])
		}

		doc := map[string]interface{}{
			"Title":       strings.Join(a.Headwords, ", "),
			"Headword":    a.Headwords,
			"HeadwordAlt": a.HeadwordsAlt,
			"Suggest":     suggests,
			"Prefix":      prefixes,
			"Content":     a.Body,
		}

		if err := json.NewEncoder(buff).Encode(map[string]interface{}{
			"index": map[string]interface{}{"_id": id},
		}); err != nil {
			return errors.Wrapf(err, "encode bulk insert meta for id %s", id)
		}

		if err := json.NewEncoder(buff).Encode(doc); err != nil {
			return errors.Wrapf(err, "encode %s doc", id)
		}

		if (i+1)%100 == 0 {
			if err := c.flushBuffer(buff); err != nil {
				return errors.Wrap(err, "flush buffer")

			}
			log.Printf("%d articles indexed", i)
			buff = &bytes.Buffer{}
		}
	}

	if err := c.flushBuffer(buff); err != nil {
		return errors.Wrap(err, "flush buffer")
	}
	log.Println("all articles indexed")

	return nil
}

func (c *commandController) assembleID(hws []string) string {
	romanized := []string{}
	for _, hw := range hws {
		// romanized = append(romanized, textutil.RomanizeBelarusian(hw))
		romanized = append(romanized, hw)
	}
	result := strings.Join(romanized, "-")
	return textutil.Slugify(result)
}

func (c *commandController) flushBuffer(buff *bytes.Buffer) error {
	err := storage.Post("/dict-"+c.indexID+"/_doc/_bulk", buff, nil)
	return errors.Wrap(err, "bulk post to storage")
}
