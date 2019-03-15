package ctl

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/verbumby/verbum/backend/pkg/storage"
	"github.com/verbumby/verbum/backend/pkg/textutil"
)

// Slugs updates slugs
func Slugs() *cobra.Command {
	return &cobra.Command{
		Use:   "slugs-update",
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

			reqbody := map[string]interface{}{
				"sort": []string{"_doc"},
				"size": 100,
			}

			type scrollbodyt struct {
				ScrollID string `json:"_scroll_id"`
				Hits     struct {
					Total int `json:"total"`
					Hits  []struct {
						Index  string `json:"_index"`
						Type   string `json:"_type"`
						ID     string `json:"_id"`
						Source struct {
							Title string
						} `json:"_source"`
					}
				} `json:"hits"`
			}
			respbody := &scrollbodyt{}
			if err := storage.Post("/dicts/_search?scroll=1m", reqbody, respbody); err != nil {
				log.Fatalf("failed to scroll over dicts: %v", err)
			}
			for len(respbody.Hits.Hits) > 0 {
				for _, hit := range respbody.Hits.Hits {
					slug := textutil.Slugify(textutil.RomanizeBelarusian(hit.Source.Title))

					if err := storage.Post("/dicts/_doc/"+hit.ID+"/_update", map[string]interface{}{
						"doc": map[string]interface{}{
							"Slug": slug,
						},
					}, nil); err != nil {
						log.Fatalf("failed up update record: %v", err)
					}
				}

				reqbody := map[string]interface{}{
					"scroll":    "1m",
					"scroll_id": respbody.ScrollID,
				}

				respbody = &scrollbodyt{}
				if err := storage.Post("/_search/scroll", reqbody, respbody); err != nil {
					log.Fatalf("failed to advance scroll over dicts: %v", err)
				}
			}

			if err := storage.Delete("/_search/scroll", map[string]interface{}{
				"scroll_id": respbody.ScrollID,
			}, nil); err != nil {
				log.Fatalf("failed to delete scroll id: %v", err)
			}
		},
	}
}
