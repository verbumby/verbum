package ctl

import (
	"encoding/json"
	"log"

	"github.com/pkg/errors"
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

			err = storage.Scroll("dicts", nil, func(rawhits []json.RawMessage) error {
				for _, rawhit := range rawhits {
					hit := &struct {
						ID     string `json:"_id"`
						Source struct {
							Title string
						} `json:"_source"`
					}{}
					if err := json.Unmarshal(rawhit, hit); err != nil {
						return errors.Wrap(err, "json unmarshal of a rawhit")
					}
					slug := textutil.Slugify(textutil.RomanizeBelarusian(hit.Source.Title))

					reqbody := map[string]interface{}{
						"doc": map[string]interface{}{
							"Slug": slug,
						},
					}
					if err := storage.Post("/dicts/_doc/"+hit.ID+"/_update", reqbody, nil); err != nil {
						return errors.Wrapf(err, "update record %s", hit.ID)
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