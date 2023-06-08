package ctl

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/verbumby/verbum/backend/storage"
	"github.com/verbumby/verbum/backend/textutil"
)

// MigrateSlugs updates slugs
func MigrateSlugs() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate-slugs-update",
		Short: "Updates slugs in dictionaries",
		Long:  "Updates slugs in dictionaries",
		Run: func(cmd *cobra.Command, args []string) {
			err := storage.Post("/dicts/_mapping/_doc", map[string]interface{}{
				"properties": map[string]interface{}{
					"Slug": map[string]interface{}{
						"type": "keyword",
					},
				},
			}, nil)
			if err != nil {
				log.Fatalf("failed to add slug field: %v", err)
			}

			err = storage.Scroll("dicts", nil, func(rawhits []json.RawMessage) error {
				for _, rawhit := range rawhits {
					hit := &struct {
						ID     string `json:"_id"`
						Source struct {
							Title string
						} `json:"_source"`
					}{}
					if err := json.Unmarshal(rawhit, hit); err != nil {
						return fmt.Errorf("json unmarshal of a rawhit: %w", err)
					}
					slug := textutil.Slugify(textutil.RomanizeBelarusian(hit.Source.Title))

					reqbody := map[string]interface{}{
						"doc": map[string]interface{}{
							"Slug": slug,
						},
					}
					if err := storage.Post("/dicts/_doc/"+hit.ID+"/_update", reqbody, nil); err != nil {
						return fmt.Errorf("update record %s: %w", hit.ID, err)
					}
				}
				return nil
			})
			if err != nil {
				log.Fatalf("failed to update slugs of dictionaries: %v", err)
			}
		},
	}
}
