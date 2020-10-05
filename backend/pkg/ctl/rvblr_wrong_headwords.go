package ctl

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/verbumby/verbum/backend/pkg/article"
	"github.com/verbumby/verbum/backend/pkg/storage"
)

// RvblrWrongHeadwords fix up rvblr wrong headwords
func RvblrWrongHeadwords() *cobra.Command {
	return &cobra.Command{
		Use:   "rvblr-wrong-headwords",
		Short: "fix up rvblr wrong headwords",
		Long:  "fix up rvblr wrong headwords",
		Run: func(cmd *cobra.Command, args []string) {
			err := storage.Scroll("dict-rvblr", nil, func(rawhits []json.RawMessage) error {
				for _, rawhit := range rawhits {
					hit := struct {
						ID     string          `json:"_id"`
						Source article.Article `json:"_source"`
					}{}
					if err := json.Unmarshal(rawhit, &hit); err != nil {
						return fmt.Errorf("unmarshal raw hit: %w", err)
					}

					a := hit.Source
					needsFix := false
					for _, hw := range a.Headword {
						if strings.ContainsAny(hw, " ()") {
							fmt.Println(hit.ID, fmt.Sprintf("'%s'", hw))
							needsFix = true
						}
					}

					if needsFix {
						contents := strings.SplitN(a.Content, "\n", 2)
						before := contents[0]
						re := regexp.MustCompile(`<\/?v-hw>`)
						after := re.ReplaceAllLiteralString(before, "")
						re = regexp.MustCompile(`[\p{Cyrillic}-]+`)
						after = re.ReplaceAllString(after, "<v-hw>$0</v-hw>")
						fmt.Println(before, "==>", after)

						contents[0] = after
						content := strings.Join(contents, "\n")

						hws, hwsalt, err := article.RvblrParse(content)
						if err != nil {
							return fmt.Errorf("couldnt parse %s new content: %w", hit.ID, err)
						}

						suggests := make([]article.Suggest, 0, len(hws)+len(hwsalt))
						for _, hw := range hws {
							suggests = append(suggests, article.Suggest{
								Input:  hw,
								Weight: 4,
							})
						}
						for _, hw := range hwsalt {
							suggests = append(suggests, article.Suggest{
								Input:  hw,
								Weight: 2,
							})
						}

						uri := fmt.Sprintf("/dict-rvblr/_doc/%s/_update", hit.ID)
						respbody := map[string]interface{}{}
						if err := storage.Post(uri, map[string]interface{}{
							"doc": map[string]interface{}{
								"Content":     content,
								"Headword":    hws,
								"HeadwordAlt": hwsalt,
								"Suggest":     suggests,
							},
						}, &respbody); err != nil {
							return fmt.Errorf("update fixed record %s: %w", hit.ID, err)
						}

						fmt.Println("Fixed!")
						fmt.Println()
					}
				}

				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
		},
	}
}
