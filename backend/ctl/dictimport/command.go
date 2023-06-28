package dictimport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/verbumby/verbum/backend/ctl/dictimport/dictparser"
	"github.com/verbumby/verbum/backend/ctl/dictimport/dictparser/dsl"
	"github.com/verbumby/verbum/backend/ctl/dictimport/dictparser/html"
	"github.com/verbumby/verbum/backend/ctl/dictimport/dictparser/stardict"
	"github.com/verbumby/verbum/backend/dictionary"
	"github.com/verbumby/verbum/backend/storage"
	"github.com/verbumby/verbum/backend/textutil"
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

	result.PersistentFlags().StringVar(&c.dictID, "dict-id", "", "dict id")
	result.PersistentFlags().StringVar(&c.indexID, "index-id", "", "storage index id")
	result.PersistentFlags().BoolVar(&c.dryrun, "dryrun", true, "true/false")
	result.PersistentFlags().BoolVarP(&c.verbose, "verbose", "v", false, "verbose output: true/false")

	return result
}

type commandController struct {
	dictID  string
	dict    dictionary.Dictionary
	indexID string
	dryrun  bool
	verbose bool
}

func (c *commandController) Run(cmd *cobra.Command, args []string) {
	if c.dryrun {
		log.Println("dryrun mode enabled")
	}
	if err := c.run(); err != nil {
		log.Fatal(err)
	}
}

func (c *commandController) getFilename() (string, error) {
	dir := viper.GetString("dicts.repo.path") + "/" + c.dictID

	files, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", dir, err)
	}

	for _, f := range files {
		if strings.HasPrefix(f.Name(), c.dictID+".") {
			return dir + "/" + f.Name(), nil
		}
	}

	return "", fmt.Errorf("couldn't find the file of %s dictionary", c.dictID)
}

func (c *commandController) run() error {
	c.dict = dictionary.GetByID(c.dictID)
	if c.dict == nil {
		return fmt.Errorf("unknown dict id %s", c.dictID)
	}

	var err error
	var d dictparser.Dictionary

	filename, err := c.getFilename()
	if err != nil {
		return err
	}
	format := path.Ext(filename)

	log.Println("processing ", filename)
	switch format {
	case ".dsl":
		d, err = dsl.ParseDSLFile(filename)
	case ".html":
		d, err = html.ParseFile(filename)
	case ".dict":
		fallthrough
	case ".stardict":
		d, err = stardict.LoadArticles(filename)
	default:
		err = fmt.Errorf("unsupported format %s", format)
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
	return storage.CreateDictIndex(c.indexID, maxResultWindow)
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
			j := 0
			for _, r := range phw {
				if j > 4 {
					break
				}
				prefix[fmt.Sprintf("Letter%d", j+1)] = string(r)
				j++
			}
			prefixes = append(prefixes, prefix)
		}

		id := strings.ToLower(a.Headwords[0])
		if d.IDsProvided {
			id = a.ID
		}
		var err error
		id, err = c.assembleID(id)
		if err != nil {
			return fmt.Errorf("assemble id for %v: %w", a.Headwords, err)
		}
		idcache[id]++
		if idcache[id] > 1 {
			id = fmt.Sprintf("%s-%d", id, idcache[id])
			log.Printf("adding index to id %s", id)
		}

		content := a.Body
		if c.dict.PrependContentWithTitle() {
			content = "<p><v-hw>" + a.Title + "</v-hw></p>\n" + content
		}

		reBrace := regexp.MustCompile(`\[.*?\]`)
		a.Title = reBrace.ReplaceAllString(a.Title, "")

		doc := map[string]interface{}{
			"Title":       a.Title,
			"Headword":    a.Headwords,
			"HeadwordAlt": a.HeadwordsAlt,
			"Phrases":     a.Phrases,
			"Suggest":     suggests,
			"Prefix":      prefixes,
			"Content":     content,
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

func (c *commandController) assembleID(firstHW string) (string, error) {
	hw := firstHW
	var romanized string
	switch c.dict.Slugifier() {
	case "belarusian":
		romanized = textutil.RomanizeBelarusian(hw)
	case "russian":
		romanized = textutil.RomanizeRussian(hw)
	case "polish":
		romanized = textutil.SlugifyPolish(hw)
	case "german":
		romanized = textutil.SlugifyDeutsch(hw)
	case "":
		romanized = hw
	default:
		return "", fmt.Errorf("unknown romanizing strategy: %s", c.dict.Slugifier())
	}
	result := romanized
	return textutil.Slugify(result), nil
}

func (c *commandController) flushBuffer(buff *bytes.Buffer) error {
	if c.dryrun {
		return nil
	}

	var resp storage.BulkResponse
	if err := storage.Post("/dict-"+c.indexID+"/_bulk", buff, &resp); err != nil {
		return fmt.Errorf("bulk post to storage: %w", err)
	}
	return resp.Error()
}
