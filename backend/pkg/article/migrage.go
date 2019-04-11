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

	if _, ok := respbody.DictRvblr.Mappings.Doc.Properties["Title"]; !ok {
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
	}

	a, err := Get("rvblr", "snieh")
	if err != nil {
		return errors.Wrap(err, "fix rvblr/snieh record: get")
	}

	if strings.HasPrefix(a.Title, " ") {
		fixedTitle := strings.TrimSpace(a.Title)
		err := storage.Post("/dict-rvblr/_doc/snieh/_update", map[string]interface{}{
			"doc": map[string]interface{}{
				"Title": fixedTitle,
			},
		}, nil)
		if err != nil {
			return errors.Wrap(err, "fix rvblr/snieh record: post")
		}
	}

	// if _, ok := respbody.DictRvblr.Mappings.Doc.Properties["Prefix1"]; !ok {
	// 	err := storage.Post("/dict-rvblr/_mapping/_doc", map[string]interface{}{
	// 		"properties": map[string]interface{}{
	// 			"Prefix1": map[string]interface{}{"type": "keyword"},
	// 			"Prefix2": map[string]interface{}{"type": "keyword"},
	// 			"Prefix3": map[string]interface{}{"type": "keyword"},
	// 		},
	// 	}, nil)
	// 	if err != nil {
	// 		return errors.Wrap(err, "add Prefix1-3 fields to dict-rvblr")
	// 	}

	// n := 0
	// err = storage.Scroll("dict-rvblr", nil, func(rawhits []json.RawMessage) error {
	// buf := &bytes.Buffer{}
	// for _, rawhit := range rawhits {
	// 	n++
	// 	hit := struct {
	// 		ID     string  `json:"_id"`
	// 		Source Article `json:"_source"`
	// 	}{}
	// 	if err := json.Unmarshal(rawhit, &hit); err != nil {
	// 		return errors.Wrap(err, "unmarshal raw hit")
	// 	}

	// 	article := hit.Source
	// 	for _, hw := range article.Headword {
	// 		if strings.ContainsAny(hw, " ()") {
	// 			fmt.Println(hit.ID, fmt.Sprintf("'%s'", hw))
	// 		}
	// 		rhw := []rune(hw)
	// 		p := Prefix{}
	// 		if len(rhw) > 0 {
	// 			p.Letter1 = string(rhw[0])
	// 		}
	// 		if len(rhw) > 1 {
	// 			p.Letter2 = string(rhw[1])
	// 		}
	// 		if len(rhw) > 2 {
	// 			p.Letter3 = string(rhw[2])
	// 		}
	// 		article.Prefix = append(article.Prefix, p)

	// 	}

	// 	for _, hw := range article.HeadwordAlt {
	// 		if strings.ContainsAny(hw, " ()") {
	// 			fmt.Println(hit.ID, hw, "HWA!")
	// 		}
	// 	}
	// fmt.Println(article.Title, "-", article.Prefix)

	// buf.WriteString(fmt.Sprintf(`{"index":{"_index":"dict-rvblr", "_type":"_doc", "_id":"%s"}}`, id))
	// buf.WriteString("\n")
	// if err := json.NewEncoder(buf).Encode(article); err != nil {
	// 	return errors.Wrap(err, "encode article")
	// }

	// fmt.Print(".")
	// }

	// respbody := struct {
	// 	Errors bool            `json:"errors"`
	// 	Items  json.RawMessage `json:"items"`
	// }{}
	// if err := storage.Post("/_bulk", buf, &respbody); err != nil {
	// 	return errors.Wrap(err, "bulk")
	// }
	// if respbody.Errors {
	// 	return fmt.Errorf("some error in one of bulk action: %s", respbody.Items)
	// }
	// fmt.Printf(" ok %d\n", n)

	// 	return nil
	// })
	// if err != nil {
	// 	return errors.Wrap(err, "migrate dict-rvblr")
	// }
	// fmt.Println("migrated", n, "records")
	// }

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
