package ctl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/verbumby/verbum/backend/pkg/article"
	"github.com/verbumby/verbum/backend/pkg/storage"
)

func MigrateStardictFixPrefixCase() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate-stardict-prefix-fix",
		Short: "Fix prefix case of stardict dictionaries",
		Long:  "Fix prefix case of stardict dictionaries",
		Run: func(cmd *cobra.Command, args []string) {
			var headwordRe = regexp.MustCompile(`(?m)<k>(.*)<\/k>`)
			n := 0

			for _, index := range []string{"dict-bel-rus", "dict-rus-bel"} {
				err := storage.Scroll(index, nil, func(rawhits []json.RawMessage) error {
					buff := &bytes.Buffer{}
					for _, h := range rawhits {
						a := struct {
							ID     string          `json:"_id"`
							Source article.Article `json:"_source"`
						}{}
						if err := json.Unmarshal(h, &a); err != nil {
							return fmt.Errorf("unmarshal %s article json: %w", h, err)
						}

						if !headwordRe.MatchString(a.Source.Content) {
							return fmt.Errorf("invalid article %s content %s", a.ID, a.Source.Content)
						}

						titleRaw := headwordRe.FindStringSubmatch(a.Source.Content)[1]
						titles := strings.Split(titleRaw, ", ")

						prefixes := []map[string]string{}
						for _, title := range titles {
							prefix := map[string]string{}
							i := 0
							for _, r := range title {
								if i > 2 {
									break
								}
								prefix[fmt.Sprintf("Letter%d", i+1)] = string(r)
								i++
							}
							prefixes = append(prefixes, prefix)
						}

						prefixesRaw, err := json.Marshal(map[string]interface{}{
							"doc": map[string]interface{}{
								"Prefix": prefixes,
							},
						})
						if err != nil {
							return fmt.Errorf("marshal prefixes of %s article: %w", a.ID, err)
						}
						// fmt.Println(a.ID, string(prefixesRaw))

						buff.WriteString(`{ "update": { "_id": "` + a.ID + `" } }`)
						buff.WriteString("\n")
						buff.Write(prefixesRaw)
						buff.WriteString("\n")
						n++
						fmt.Print(".")
					}

					if err := storage.Post("/"+index+"/_bulk", buff, nil); err != nil {
						return fmt.Errorf("bulk update: %w", err)
					}

					buff = &bytes.Buffer{}
					fmt.Println(" ", n)
					return nil
				})
				if err != nil {
					log.Fatal(err)
				}
			}
		},
	}
}
