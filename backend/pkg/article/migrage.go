package article

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/verbumby/verbum/backend/pkg/dictionary"
	"github.com/verbumby/verbum/backend/pkg/storage"
	"github.com/verbumby/verbum/backend/pkg/textutil"
)

// Migrate migrates article storage
func Migrate() error {
	var rvblrMapping map[string]interface{}
	{
		respbody := struct {
			DictRvblr struct {
				Mappings struct {
					Properties map[string]interface{} `json:"properties"`
				} `json:"mappings"`
			} `json:"dict-rvblr"`
		}{}
		if err := storage.Get("/dict-rvblr/_mappings", &respbody); err != nil {
			return fmt.Errorf("get dict-rvblr mappings: %w", err)
		}
		rvblrMapping = respbody.DictRvblr.Mappings.Properties
	}

	type indexSettingsType struct {
		MaxResultWindow string `json:"max_result_window"`
	}
	var rvblrSettings indexSettingsType
	{
		respbody := struct {
			DictRvblr struct {
				Settings struct {
					Index indexSettingsType `json:"index"`
				} `json:"settings"`
			} `json:"dict-rvblr"`
		}{}
		if err := storage.Get("/dict-rvblr/_settings", &respbody); err != nil {
			return fmt.Errorf("get dict-rvblr settings: %w", err)
		}
		rvblrSettings = respbody.DictRvblr.Settings.Index
	}

	rvblrAddTitleMigration := func() error {
		if _, ok := rvblrMapping["Title"]; ok {
			return nil
		}
		err := storage.Post("/dict-rvblr/_mapping/_doc", map[string]interface{}{
			"properties": map[string]interface{}{
				"Title": map[string]interface{}{
					"type": "keyword",
				},
			},
		}, nil)
		if err != nil {
			return fmt.Errorf("add Title field to dict-rvblr: %w", err)
		}

		re := regexp.MustCompile(`<\/?v-hw>`)
		slugmap := map[string]int{}
		n := 0
		err = storage.Scroll("dict-rvblr", nil, func(rawhits []json.RawMessage) error {
			buf := &bytes.Buffer{}
			for _, rawhit := range rawhits {
				n++
				hit := struct {
					ID     string  `json:"_id"`
					Source Article `json:"_source"`
				}{}
				if err := json.Unmarshal(rawhit, &hit); err != nil {
					return fmt.Errorf("unmarshal raw hit: %w", err)
				}

				article := hit.Source
				parts := strings.SplitN(article.Content, "\n", 2)
				article.Title = re.ReplaceAllString(parts[0], "")
				slug := textutil.Slugify(textutil.RomanizeBelarusian(article.Title))
				slugmap[slug]++
				id := slug
				if slugmap[slug] > 1 {
					id += fmt.Sprintf("-%d", slugmap[slug])
				}

				buf.WriteString(fmt.Sprintf(`{"index":{"_index":"dict-rvblr", "_type":"_doc", "_id":"%s"}}`, id))
				buf.WriteString("\n")
				if err := json.NewEncoder(buf).Encode(article); err != nil {
					return fmt.Errorf("encode article: %w", err)
				}

				buf.WriteString(fmt.Sprintf(`{"delete":{"_index":"dict-rvblr", "_type":"_doc", "_id":"%s"}}`, hit.ID))
				buf.WriteString("\n")
				fmt.Print(".")
			}

			respbody := struct {
				Errors bool            `json:"errors"`
				Items  json.RawMessage `json:"items"`
			}{}
			if err := storage.Post("/_bulk", buf, &respbody); err != nil {
				return fmt.Errorf("bulk: %w", err)
			}
			if respbody.Errors {
				return fmt.Errorf("some error in one of bulk action: %s", respbody.Items)
			}
			fmt.Printf(" ok %d\n", n)

			return nil
		})
		if err != nil {
			return fmt.Errorf("migrate dict-rvblr: %w", err)
		}
		fmt.Println("migrated", n, "records")
		return nil
	}

	rvblrFixSniehTitle := func() error {
		a, err := Get(dictionary.Get("rvblr"), "snieh")
		if err != nil {
			return fmt.Errorf("get: %w", err)
		}

		if strings.HasPrefix(a.Title, " ") {
			fixedTitle := strings.TrimSpace(a.Title)
			err := storage.Post("/dict-rvblr/_doc/snieh/_update", map[string]interface{}{
				"doc": map[string]interface{}{
					"Title": fixedTitle,
				},
			}, nil)
			if err != nil {
				return fmt.Errorf("post: %w", err)
			}
		}
		return nil
	}

	rvblrPrefix := func() error {
		if _, ok := rvblrMapping["Prefix"]; ok {
			return nil
		}
		err := storage.Post("/dict-rvblr/_mapping/_doc", map[string]interface{}{
			"properties": map[string]interface{}{
				"Prefix": map[string]interface{}{
					"type": "nested",
					"properties": map[string]interface{}{
						"Letter1": map[string]interface{}{"type": "keyword"},
						"Letter2": map[string]interface{}{"type": "keyword"},
						"Letter3": map[string]interface{}{"type": "keyword"},
					},
				},
			},
		}, nil)
		if err != nil {
			return fmt.Errorf("add Prefix nested field to dict-rvblr: %w", err)
		}

		n := 0
		err = storage.Scroll("dict-rvblr", nil, func(rawhits []json.RawMessage) error {
			buf := &bytes.Buffer{}
			for _, rawhit := range rawhits {
				n++
				hit := struct {
					ID     string  `json:"_id"`
					Source Article `json:"_source"`
				}{}
				if err := json.Unmarshal(rawhit, &hit); err != nil {
					return fmt.Errorf("unmarshal raw hit: %w", err)
				}

				a := hit.Source
				seen := map[string]bool{}
				for _, hw := range a.Headword {
					rhw := []rune(hw)

					if _, ok := seen[string(rhw[:3])]; ok {
						continue
					}
					seen[string(rhw[:3])] = true

					p := Prefix{}
					if len(rhw) > 0 {
						p.Letter1 = string(rhw[0])
					}
					if len(rhw) > 1 {
						p.Letter2 = string(rhw[1])
					}
					if len(rhw) > 2 {
						p.Letter3 = string(rhw[2])
					}
					a.Prefix = append(a.Prefix, p)
				}

				buf.WriteString(fmt.Sprintf(`{"update":{"_index":"dict-rvblr", "_type":"_doc", "_id":"%s"}}`, hit.ID))
				buf.WriteString("\n")
				if err := json.NewEncoder(buf).Encode(map[string]interface{}{
					"doc": map[string]interface{}{
						"Prefix": a.Prefix,
					},
				}); err != nil {
					return fmt.Errorf("encode article: %w", err)
				}

				fmt.Print(".")
			}

			respbody := struct {
				Errors bool            `json:"errors"`
				Items  json.RawMessage `json:"items"`
			}{}
			if err := storage.Post("/_bulk", buf, &respbody); err != nil {
				return fmt.Errorf("bulk: %w", err)
			}
			if respbody.Errors {
				return fmt.Errorf("some error in one of bulk action: %s", respbody.Items)
			}
			fmt.Printf(" ok %d\n", n)

			return nil
		})
		if err != nil {
			return fmt.Errorf("migrate dict-rvblr: %w", err)
		}
		fmt.Println("migrated", n, "records")
		return nil
	}

	rvblrMaxResultWindow := func() error {
		if rvblrSettings.MaxResultWindow == "40000" {
			return nil
		}

		if err := storage.Query(http.MethodPut, "/dict-rvblr/_settings", map[string]interface{}{
			"index": map[string]interface{}{
				"max_result_window": "40000",
			},
		}, nil); err != nil {
			return fmt.Errorf("post settings: %w", err)
		}
		return nil
	}

	migrations := []struct {
		name string
		f    func() error
	}{
		{
			name: "rvblr add title field",
			f:    rvblrAddTitleMigration,
		},
		{
			name: "rvblr fix snieh record title",
			f:    rvblrFixSniehTitle,
		},
		{
			name: "rvblr add prefix fields and data",
			f:    rvblrPrefix,
		},
		{
			name: "update max result window index setting",
			f:    rvblrMaxResultWindow,
		},
	}

	for _, m := range migrations {
		if err := m.f(); err != nil {
			return fmt.Errorf("migration %s: %w", m.name, err)
		}
	}

	return nil
}
