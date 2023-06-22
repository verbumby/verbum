package ctl

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/verbumby/verbum/backend/article"
	"github.com/verbumby/verbum/backend/dictionary"
	"github.com/verbumby/verbum/backend/storage"
)

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

	sort.Slice(hits, func(i, j int) bool {
		return hits[i].ID < hits[j].ID
	})

	for _, a := range hits {
		r := a.Source.Dictionary.ToHTML(a.Source.Content)
		rs := string(r)

		if first {
			first = false
		} else {
			fmt.Println("<hr/>")
		}

		if strings.HasSuffix(rs, "\n") {
			fmt.Print(rs)
		} else {
			fmt.Println(rs)
		}
	}

	return nil
}
