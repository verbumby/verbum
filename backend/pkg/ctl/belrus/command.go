package belrus

import (
	"bytes"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/verbumby/verbum/backend/pkg/storage"
	"github.com/verbumby/verbum/backend/pkg/textutil"
	"golang.org/x/text/unicode/norm"
)

// Command creates cobra command
func Command() *cobra.Command {
	var filename string
	result := &cobra.Command{
		Use:   "bel-rus-import",
		Short: "Imports bel-rus dictionary",
		Long:  "Imports bel-rus dictionary",
		Run: func(cmd *cobra.Command, args []string) {
			as, err := parseFile(filename)
			if err != nil {
				log.Fatalf("parser failed: %v", err)
			}
			fmt.Printf("parsed %d articles\n", len(as))

			if err := createIndex(); err != nil {
				log.Fatalf("create index: %v", err)
			}
			log.Println("created index")

			if err := indexArticles(as); err != nil {
				log.Fatalf("index articles: %v", err)
			}
		},
	}

	result.PersistentFlags().StringVar(&filename, "filename", "", "filename of the dict")
	return result
}

func createIndex() error {
	respbody := map[string]interface{}{}
	err := storage.Put("/dict-bel-rus", map[string]interface{}{
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

func parseFile(filename string) ([]parsedArticle, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "os open file")
	}
	defer f.Close()

	c, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Wrap(err, "read all")
	}

	n := 0
	result := []parsedArticle{}
	for {
		n++
		i := bytes.Index(c[3:], []byte("<k>"))
		if i == -1 {
			result = append(result, processArticle(c))
			break
		}
		i += 3
		result = append(result, processArticle(c[:i]))
		c = c[i:]
	}
	return result, nil
}

type parsedArticle struct {
	title    string
	slug     string
	hws      []string
	content  string
	prefixes []map[string]string
	suggests []map[string]interface{}
}

var re = regexp.MustCompile(`<k>(.*)<\/k>`)

func processArticle(s []byte) parsedArticle {
	ss := html.UnescapeString(string(s))
	ss = norm.NFC.String(ss)

	sms := re.FindAllStringSubmatch(ss, -1)
	if len(sms) != 1 {
		panic("invalid count of headwords in: " + ss)
	}

	title := sms[0][1]
	slug := textutil.Slugify(textutil.RomanizeBelarusian(title))

	hws := []string{}
	for _, hw := range strings.Split(title, ",") {
		hws = append(hws, strings.TrimSpace(hw))
	}

	suggests := []map[string]interface{}{}
	for _, hw := range hws {
		hw := strings.ToLower(hw)
		suggests = append(suggests, map[string]interface{}{
			"input":  hw,
			"weight": 4,
		})
	}

	prefixes := []map[string]string{}
	for _, hw := range hws {
		hw := strings.ToLower(hw)
		prefix := map[string]string{}
		i := 0
		for _, r := range hw {
			if i > 2 {
				break
			}
			prefix[fmt.Sprintf("Letter%d", i+1)] = string(r)
			i++
		}
		prefixes = append(prefixes, prefix)
	}

	return parsedArticle{
		title:    title,
		content:  ss,
		hws:      hws,
		slug:     slug,
		prefixes: prefixes,
		suggests: suggests,
	}
}

func indexArticles(as []parsedArticle) error {
	idcache := map[string]int{}
	for i, a := range as {
		idcache[a.slug]++
		id := a.slug
		if idcache[a.slug] > 1 {
			id = fmt.Sprintf("%s-%d", a.slug, idcache[a.slug])
		}
		if err := storage.Put("/dict-bel-rus/_doc/"+id, map[string]interface{}{
			"Title":    a.title,
			"Headword": a.hws,
			"Suggest":  a.suggests,
			"Prefix":   a.prefixes,
			"Content":  a.content,
		}, nil); err != nil {
			return errors.Wrapf(err, "storage put of %s", a.slug)
		}
		fmt.Print(".")
		if i%120 == 0 {
			fmt.Println(" ", i)
		}
	}

	return nil
}
