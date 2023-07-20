package dictimport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/verbumby/verbum/backend/ctl/dictimport/dictparser"
	"github.com/verbumby/verbum/backend/ctl/dictimport/dictparser/dsl"
	"github.com/verbumby/verbum/backend/ctl/dictimport/dictparser/grammardb"
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
	dictID     string
	dict       dictionary.Dictionary
	indexID    string
	dryrun     bool
	verbose    bool
	useDictIDs bool
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

	if c.dictID == "grammardb" {
		return dir, nil
	}

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

var reIndexSuffix = regexp.MustCompile(`^dict-(?:.+?)(?:-(\d*))?$`)

func (c *commandController) autoIndexID() (string, error) {
	path := fmt.Sprintf("/dict-%s-*", c.dictID)
	resp := map[string]any{}
	if err := storage.Get(path, &resp); err != nil {
		return "", fmt.Errorf("list indices: %w", err)
	}
	var d int64 = 0

	for index := range resp {
		m := reIndexSuffix.FindStringSubmatch(index)
		if m == nil {
			continue
		}

		if m[1] == "" {
			continue
		}

		nd, err := strconv.ParseInt(m[1], 10, 32)
		if err != nil {
			return "", fmt.Errorf("can't parse %s as int: %w", m[1], err)
		}

		if nd > d {
			d = nd
		}
	}

	return fmt.Sprintf("%s-%d", c.dictID, d+1), nil
}

func (c *commandController) run() error {
	var err error

	c.dict = dictionary.GetByID(c.dictID)
	if c.dict == nil {
		return fmt.Errorf("unknown dict id %s", c.dictID)
	}

	if c.indexID == "" {
		c.indexID, err = c.autoIndexID()
		if err != nil {
			return fmt.Errorf("calculate index name: %w", err)
		}
	}
	log.Printf("indexing into %s", c.indexID)

	filename, err := c.getFilename()
	if err != nil {
		return err
	}
	log.Println("processing ", filename)

	var articlesCh chan dictparser.Article
	var errCh chan error

	switch c.dict.(type) {
	case dictionary.GrammarDB:
		articlesCh, errCh = grammardb.ParseDirectory(filename)
		c.useDictIDs = true

	case dictionary.DSL:
		file, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("open %s", filename)
		}
		defer file.Close()

		articlesCh, errCh = dsl.ParseReader(file)

	case dictionary.HTML:
		file, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("open %s", filename)
		}
		defer file.Close()

		c.useDictIDs = true
		articlesCh, errCh = html.ParseReader(file)

	case dictionary.Stardict:
		file, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("open %s", filename)
		}
		defer file.Close()

		articlesCh, errCh = stardict.LoadArticles(file)

	default:
		err = fmt.Errorf("unsupported format %T", c.dict)
	}
	if err != nil {
		return fmt.Errorf("parse dictionary: %w", err)
	}

	if err := c.createIndex(); err != nil {
		return fmt.Errorf("create index: %w", err)
	}

	if err := c.indexArticles(articlesCh); err != nil {
		return fmt.Errorf("index articles: %w", err)
	}

	if err := <-errCh; err != nil {
		return err
	}

	log.Println("all articles indexed")

	if err := c.updateAlias(); err != nil {
		return fmt.Errorf("update alias: %w", err)
	}

	return nil
}

func (c *commandController) createIndex() error {
	if c.dryrun {
		return nil
	}
	return storage.CreateDictIndex(c.indexID)
}

func (c *commandController) indexArticles(articlesCh chan dictparser.Article) error {
	idcache := map[string]int{}
	iddups := map[string]int{}

	buff := &bytes.Buffer{}
	i := -1
	for a := range articlesCh {
		i++
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
		if c.useDictIDs {
			id = a.ID
			iddups[id]++
			if iddups[id] > 1 {
				return fmt.Errorf("duplicate id: %s", id)
			}
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

	return nil
}

func (c *commandController) assembleID(firstHW string) (string, error) {
	hw := firstHW
	var romanized string
	switch c.dict.Slugifier() {
	case "belarusian":
		romanized = textutil.RomanizeBelarusian(hw)
	case "none":
		return firstHW, nil
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

func (c *commandController) updateAlias() error {
	alias := "dict-" + c.dictID
	path := fmt.Sprintf("/dict-%s-*", c.dictID)
	resp := map[string]struct {
		Aliases map[string]any `json:"aliases"`
	}{}
	if err := storage.Get(path, &resp); err != nil {
		return fmt.Errorf("list indexes: %w", err)
	}

	toBeRemoved := []string{}

	for indexName, index := range resp {
		if _, ok := index.Aliases[alias]; ok {
			toBeRemoved = append(toBeRemoved, indexName)
		}
	}

	toBeAdded := "dict-" + c.indexID

	log.Printf("removing %v and adding %s to %s alias", toBeRemoved, toBeAdded, alias)

	if c.dryrun {
		return nil
	}

	actions := []any{}
	for _, index := range toBeRemoved {
		actionBody := map[string]any{
			"index": index,
			"alias": alias,
		}

		action := map[string]any{
			"remove": actionBody,
		}
		actions = append(actions, action)
	}

	{
		actionBody := map[string]any{
			"index": toBeAdded,
			"alias": alias,
		}

		action := map[string]any{
			"add": actionBody,
		}
		actions = append(actions, action)
	}

	if err := storage.Post("/_aliases", map[string]any{"actions": actions}, nil); err != nil {
		return fmt.Errorf("update alias %s by removing %v and adding %s: %w", alias, toBeRemoved, toBeAdded, err)
	}

	return nil
}
