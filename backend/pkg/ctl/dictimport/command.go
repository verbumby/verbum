package dictimport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/verbumby/verbum/backend/pkg/ctl/dictimport/dictparser"
	"github.com/verbumby/verbum/backend/pkg/ctl/dictimport/dictparser/dsl"
	"github.com/verbumby/verbum/backend/pkg/ctl/dictimport/dictparser/html"
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

	result.PersistentFlags().StringVar(&c.filename, "filename", "", "filename/dictionary of the dict")
	result.PersistentFlags().StringVar(&c.format, "format", "", "dsl|stardict|html")
	result.PersistentFlags().StringVar(&c.indexID, "index-id", "", "storage index id")
	result.PersistentFlags().StringVar(&c.romanizer, "romanizer", "", "<blank>|belarusian|russian")
	result.PersistentFlags().BoolVar(&c.dryrun, "dryrun", true, "true/false")
	result.PersistentFlags().IntVar(&c.limit, "limit", 1000, "limits the number of articles processed, -1 disables the limit")
	result.PersistentFlags().BoolVarP(&c.verbose, "verbose", "v", false, "verbose output: true/false")

	return result
}

type commandController struct {
	filename  string
	format    string
	indexID   string
	romanizer string
	dryrun    bool
	limit     int
	verbose   bool
}

func (c *commandController) Run(cmd *cobra.Command, args []string) {
	if !c.dryrun || c.limit == -1 {
		c.limit = math.MaxInt32
	}
	if c.dryrun {
		log.Println("dryrun mode enabled")
	}
	log.Println("processing ", c.filename)
	if err := c.run(); err != nil {
		log.Fatal(err)
	}
}

func (c *commandController) run() error {
	var err error
	var d dictparser.Dictionary
	switch c.format {
	case "dsl":
		d, err = dsl.ParseDSLFile(c.filename)
	case "html":
		d, err = html.LoadArticles(c.filename)
	default:
		err = fmt.Errorf("unsupported format %s", c.format)
	}
	if err != nil {
		return fmt.Errorf("parse dictionary: %w", err)
	}

	log.Printf("found %d articles in the dictionary", len(d.Articles))

	if err := c.createIndex(len(d.Articles) + 50000); err != nil {
		return fmt.Errorf("create index: %w", err)
	}

	if err := c.indexArticles(d); err != nil {
		return fmt.Errorf("index articles: %w", err)
	}

	return nil
}

func (c *commandController) createIndex(maxResultWindow int) error {
	if c.dryrun {
		return nil
	}
	err := storage.Put("/dict-"+c.indexID, map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 0,
			"max_result_window":  maxResultWindow,
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
						"tokenize_on_chars": []string{"-", ".", "/", "—", " ", "(", ")", ",", "!", "?", "…"},
					},
				},
			},
		},
		"mappings": map[string]interface{}{
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
				"Phrases": map[string]any{
					"type":     "text",
					"analyzer": "standard",
				},
				"Prefix": map[string]interface{}{
					"type": "nested",
					"properties": map[string]interface{}{
						"Letter1": map[string]interface{}{"type": "keyword"},
						"Letter2": map[string]interface{}{"type": "keyword"},
						"Letter3": map[string]interface{}{"type": "keyword"},
						"Letter4": map[string]interface{}{"type": "keyword"},
						"Letter5": map[string]interface{}{"type": "keyword"},
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
				"ModifiedAt": map[string]interface{}{
					"type": "date",
				},
			},
		},
	}, nil)
	if err != nil {
		return fmt.Errorf("storage put: %w", err)
	}
	return nil
}

func (c *commandController) indexArticles(d dictparser.Dictionary) error {
	idcache := map[string]int{}

	buff := &bytes.Buffer{}
	for i, a := range d.Articles {
		if i == c.limit {
			break
		}

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
				if i > 4 {
					break
				}
				prefix[fmt.Sprintf("Letter%d", i+1)] = string(r)
				i++
			}
			prefixes = append(prefixes, prefix)
		}

		var id string
		if d.IDsProvided {
			id = textutil.Slugify(a.ID)
		} else {
			var err error
			id, err = c.assembleID(a.Headwords)
			if err != nil {
				return fmt.Errorf("assemble id for %v: %w", a.Headwords, err)
			}
			idcache[id]++
			if idcache[id] > 1 {
				id = fmt.Sprintf("%s-%d", id, idcache[id])
			}
		}

		doc := map[string]interface{}{
			"Title":       strings.Join(a.Headwords, ", "),
			"Headword":    a.Headwords,
			"HeadwordAlt": a.HeadwordsAlt,
			"Phrases":     a.Phrases,
			"Suggest":     suggests,
			"Prefix":      prefixes,
			"Content":     a.Body,
			"ModifiedAt":  time.Now().UTC().Format(time.RFC3339),
		}

		if err := json.NewEncoder(buff).Encode(map[string]interface{}{
			"create": map[string]any{"_id": id},
		}); err != nil {
			return fmt.Errorf("encode bulk insert meta for id %s: %w", id, err)
		}

		if err := json.NewEncoder(buff).Encode(doc); err != nil {
			return fmt.Errorf("encode %s doc: %w", id, err)
		}

		if c.verbose {
			toprint := map[string]interface{}{"_doc": doc, "_id": id}
			if err := json.NewEncoder(os.Stdout).Encode(toprint); err != nil {
				return fmt.Errorf("encode %s doc for verbose output: %w", id, err)
			}
			fmt.Println()
		}

		if (i+1)%100 == 0 {
			if err := c.flushBuffer(buff); err != nil {
				return fmt.Errorf("flush buffer: %w", err)

			}
			log.Printf("%d articles indexed", i)
			buff = &bytes.Buffer{}
		}
	}

	if err := c.flushBuffer(buff); err != nil {
		return fmt.Errorf("flush buffer: %w", err)
	}
	log.Println("all articles indexed")

	return nil
}

func (c *commandController) assembleID(hws []string) (string, error) {
	romanized := []string{}
	for _, hw := range hws {
		switch c.romanizer {
		case "belarusian":
			romanized = append(romanized, textutil.RomanizeBelarusian(hw))
		case "russian":
			romanized = append(romanized, textutil.RomanizeRussian(hw))
		case "":
			romanized = append(romanized, hw)
		default:
			return "", fmt.Errorf("unknown romanizing strategy: %s", c.romanizer)
		}
	}
	result := strings.Join(romanized, "-")
	return textutil.Slugify(result), nil
}

func (c *commandController) flushBuffer(buff *bytes.Buffer) error {
	if c.dryrun {
		return nil
	}
	type respItemType struct {
		ID    string          `json:"_id"`
		Error json.RawMessage `json:"error"`
	}

	type respType struct {
		Errors bool `json:"errors"`
		Items  []struct {
			Create *respItemType `json:"create"`
			Index  *respItemType `json:"index"`
			Delete *respItemType `json:"delete"`
			Update *respItemType `json:"update"`
		} `json:"items"`
	}

	var resp respType
	if err := storage.Post("/dict-"+c.indexID+"/_doc/_bulk", buff, &resp); err != nil {
		return fmt.Errorf("bulk post to storage: %w", err)
	}

	if resp.Errors {
		errors := []string{}
		for _, item := range resp.Items {
			if item.Create.Error != nil {
				errors = append(errors, fmt.Sprintf("id `%s`: %s", item.Create.ID, string(item.Create.Error)))
			}
		}
		return fmt.Errorf("bulk post to storage returned errors: " + strings.Join(errors, "; "))
	}
	return nil
}
