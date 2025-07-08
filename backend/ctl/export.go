package ctl

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/verbumby/verbum/backend/article"
	"github.com/verbumby/verbum/backend/dictionary"
	"github.com/verbumby/verbum/backend/storage"
)

var reStrongSpanDarkBlue = regexp.MustCompile(`<strong><span style="color: darkblue">([^<]*)</span></strong>`)

func ExportCommand() *cobra.Command {
	c := &exportController{}
	result := &cobra.Command{
		Use: "export",
		Run: c.Run,
	}

	result.PersistentFlags().StringVar(&c.indexID, "index-id", "", "Source dictionary index ID")
	result.PersistentFlags().StringVar(&c.format, "format", "", "Format of the exported file. Only HTML is currently supported")

	return result
}

type exportController struct {
	indexID string
	format  string
}

func (c *exportController) Run(cmd *cobra.Command, args []string) {
	if err := c.run(); err != nil {
		log.Fatal(err)
	}
}

func (c *exportController) run() error {
	first := true
	type hit struct {
		Source article.Article `json:"_source"`
		Index  string          `json:"_index"`
		ID     string          `json:"_id"`
	}

	hits := []hit{}

	err := storage.Scroll("dict-"+c.indexID, nil, func(rawhits []json.RawMessage) error {
		for _, rawhit := range rawhits {
			a := hit{}
			if err := json.Unmarshal(rawhit, &a); err != nil {
				return fmt.Errorf("unmarshal %s article json: %w", rawhit, err)
			}
			a.Source.Dictionary = dictionary.GetByIndex(a.Index)
			hits = append(hits, a)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("scroll through the source index: %w", err)
	}

	for _, a := range hits {
		content := a.Source.Content
		// content = strings.ReplaceAll(content, "<", "&lt;")
		// content = strings.ReplaceAll(content, ">", "&gt;")
		r := a.Source.Dictionary.ToHTML(content)
		rs := string(r)

		if first {
			first = false
		} else {
			fmt.Println("<hr/>")
		}

		rs = reStrongSpanDarkBlue.ReplaceAllString(rs, `<strong class="hw">$1</strong>`)

		rs = strings.Replace(rs, `<strong class="hw">`, `<strong class="hw" id="`+a.ID+`">`, 1)
		rs = strings.ReplaceAll(rs, `</p><p`, "</p>\n<p")

		if strings.HasSuffix(rs, "\n") {
			fmt.Print(rs)
		} else {
			fmt.Println(rs)
		}
	}

	return nil
}
