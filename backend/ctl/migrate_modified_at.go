package ctl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/verbumby/verbum/backend/dictionary"
	"github.com/verbumby/verbum/backend/storage"
)

func MigrateModifiedAt() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate-modified-at",
		Short: "Adds current timestamp as ModifiedAt to all articles",
		Long:  "Adds current timestamp as ModifiedAt to all articles",
		RunE: func(cmd *cobra.Command, args []string) error {
			dicts := dictionary.GetAll()

			for _, d := range dicts {
				log.Printf("processing %s dict", d.ID())
				n := 0
				index := "dict-" + d.ID()
				err := storage.Scroll(index, nil, func(rawhits []json.RawMessage) error {
					buff := &bytes.Buffer{}
					for _, rawhit := range rawhits {
						hit := struct {
							ID     string `json:"_id"`
							Source struct {
								ModifiedAt string `json:"ModifiedAt"`
							} `json:"_source"`
						}{}
						if err := json.Unmarshal(rawhit, &hit); err != nil {
							return fmt.Errorf("unmarshal %s hit: %w", string(rawhit), err)
						}
						if hit.Source.ModifiedAt != "" {
							continue
						}

						modifiedAt, err := json.Marshal(map[string]interface{}{
							"doc": map[string]interface{}{
								"ModifiedAt": time.Now().UTC().Format(time.RFC3339),
							},
						})
						if err != nil {
							return fmt.Errorf("marshal prefixes of %s article: %w", hit.ID, err)
						}

						buff.WriteString(`{ "update": { "_id": "` + hit.ID + `" } }`)
						buff.WriteString("\n")
						buff.Write(modifiedAt)
						buff.WriteString("\n")

						n++
						fmt.Print(".")
					}

					if buff.Len() == 0 {
						return nil
					}

					if err := storage.Post("/"+index+"/_bulk", buff, nil); err != nil {
						return fmt.Errorf("bulk update: %w", err)
					}

					fmt.Println(" ", n)
					return nil
				})
				if err != nil {
					return fmt.Errorf("migrate dict %s: %w", d.ID(), err)
				}
			}
			return nil
		},
	}
}
