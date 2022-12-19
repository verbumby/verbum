package ctl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/verbumby/verbum/backend/pkg/storage"
)

func ReindexCommand() *cobra.Command {
	c := &reindexController{}
	result := &cobra.Command{
		Use: "reindex",
		Run: c.Run,
	}

	result.PersistentFlags().StringVar(&c.sourceIndexID, "source-index-id", "", "Source dictionary index ID")
	result.PersistentFlags().StringVar(&c.targetIndexID, "target-index-id", "", "Target dictionary index ID")

	return result
}

type reindexController struct {
	sourceIndexID string
	targetIndexID string
}

func (c *reindexController) Run(cmd *cobra.Command, args []string) {
	if err := c.run(); err != nil {
		log.Fatal(err)
	}
}

func (c *reindexController) run() error {
	sourceIndexSettings := map[string]struct {
		Settings struct {
			Index struct {
				MaxResultWindow string `json:"max_result_window"`
			} `json:"index"`
		} `json:"settings"`
	}{}
	if err := storage.Get("/dict-"+c.sourceIndexID+"/_settings/index.max_result_window", &sourceIndexSettings); err != nil {
		return fmt.Errorf("get source index max result window: %w", err)
	}
	maxResultWindow, err := strconv.Atoi(sourceIndexSettings["dict-"+c.sourceIndexID].Settings.Index.MaxResultWindow)
	if err != nil {
		return fmt.Errorf("convert max result window value to int: %w", err)
	}

	if err := storage.CreateDictIndex(c.targetIndexID, maxResultWindow); err != nil {
		return fmt.Errorf("create target index: %w", err)
	}

	count := 0
	buff := &bytes.Buffer{}
	err = storage.Scroll("dict-"+c.sourceIndexID, nil, func(rawhits []json.RawMessage) error {
		for _, rawhit := range rawhits {
			a := struct {
				ID     string         `json:"_id"`
				Source map[string]any `json:"_source"`
			}{}
			if err := json.Unmarshal(rawhit, &a); err != nil {
				return fmt.Errorf("unmarshal %s article json: %w", rawhit, err)
			}

			if err := json.NewEncoder(buff).Encode(map[string]any{
				"create": map[string]any{"_id": a.ID},
			}); err != nil {
				return fmt.Errorf("encode bulk insert meta for id %s: %w", a.ID, err)
			}

			if err := json.NewEncoder(buff).Encode(a.Source); err != nil {
				return fmt.Errorf("encode %s doc: %w", a.ID, err)
			}
			count++
		}

		if err := c.flushBuffer(buff); err != nil {
			return fmt.Errorf("flush buffer: %w", err)
		}

		buff.Reset()

		log.Printf("Reindexed %d articles", count)
		return nil
	})
	if err != nil {
		return fmt.Errorf("scroll through the source index: %w", err)
	}

	return nil
}

func (c *reindexController) flushBuffer(buff *bytes.Buffer) error {
	var resp storage.BulkResponse
	if err := storage.Post("/dict-"+c.targetIndexID+"/_doc/_bulk", buff, &resp); err != nil {
		return fmt.Errorf("bulk post to storage: %w", err)
	}
	return resp.Error()
}
