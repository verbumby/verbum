package article

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
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
					Doc struct {
						Properties map[string]interface{} `json:"properties"`
					} `json:"_doc"`
				} `json:"mappings"`
			} `json:"dict-rvblr"`
		}{}
		if err := storage.Get("/dict-rvblr/_mappings", &respbody); err != nil {
			return errors.Wrap(err, "get dict-rvblr mappings")
		}
		rvblrMapping = respbody.DictRvblr.Mappings.Doc.Properties
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
			return errors.Wrap(err, "add Title field to dict-rvblr")
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
					return errors.Wrap(err, "unmarshal raw hit")
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
					return errors.Wrap(err, "encode article")
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
				return errors.Wrap(err, "bulk")
			}
			if respbody.Errors {
				return fmt.Errorf("some error in one of bulk action: %s", respbody.Items)
			}
			fmt.Printf(" ok %d\n", n)

			return nil
		})
		if err != nil {
			return errors.Wrap(err, "migrate dict-rvblr")
		}
		fmt.Println("migrated", n, "records")
		return nil
	}

	rvblrFixSniehTitle := func() error {
		a, err := Get("rvblr", "snieh")
		if err != nil {
			return errors.Wrap(err, "get")
		}

		if strings.HasPrefix(a.Title, " ") {
			fixedTitle := strings.TrimSpace(a.Title)
			err := storage.Post("/dict-rvblr/_doc/snieh/_update", map[string]interface{}{
				"doc": map[string]interface{}{
					"Title": fixedTitle,
				},
			}, nil)
			if err != nil {
				return errors.Wrap(err, "post")
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
			return errors.Wrap(err, "add Prefix nested field to dict-rvblr")
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
					return errors.Wrap(err, "unmarshal raw hit")
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
					return errors.Wrap(err, "encode article")
				}

				fmt.Print(".")
			}

			respbody := struct {
				Errors bool            `json:"errors"`
				Items  json.RawMessage `json:"items"`
			}{}
			if err := storage.Post("/_bulk", buf, &respbody); err != nil {
				return errors.Wrap(err, "bulk")
			}
			if respbody.Errors {
				return fmt.Errorf("some error in one of bulk action: %s", respbody.Items)
			}
			fmt.Printf(" ok %d\n", n)

			return nil
		})
		if err != nil {
			return errors.Wrap(err, "migrate dict-rvblr")
		}
		fmt.Println("migrated", n, "records")
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
	}

	for _, m := range migrations {
		if err := m.f(); err != nil {
			return errors.Wrapf(err, m.name)
		}
	}

	return nil
}

func stringArrayContains(a []string, s string) bool {
	for _, as := range a {
		if as == s {
			return true
		}
	}
	return false
}
