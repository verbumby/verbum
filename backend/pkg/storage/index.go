package storage

import "fmt"

func CreateDictIndex(dictID string, maxResultWindow int) error {
	err := Put("/dict-"+dictID, map[string]any{
		"settings": map[string]any{
			"number_of_shards":   1,
			"number_of_replicas": 0,
			"max_result_window":  maxResultWindow,
			"analysis": map[string]any{
				"analyzer": map[string]any{
					"headword": map[string]any{
						"filter":    []string{"lowercase"},
						"type":      "custom",
						"tokenizer": "keyword",
					},
					"headword_smaller": map[string]any{
						"filter":    []string{"lowercase"},
						"type":      "custom",
						"tokenizer": "headword_smaller",
					},
					"body": map[string]any{
						"type":        "custom",
						"tokenizer":   "standard",
						"char_filter": []string{"html_strip", "dsl_strip"},
						"filter":      []string{"lowercase"},
					},
				},
				"tokenizer": map[string]any{
					"headword_smaller": map[string]any{
						"type":              "char_group",
						"tokenize_on_chars": []string{"-", ".", "/", "—", " ", "(", ")", ",", "!", "?", "…"},
					},
				},
				"char_filter": map[string]any{
					"dsl_strip": map[string]any{
						"type":    "pattern_replace",
						"pattern": "\\[/?.*?\\]",
					},
				},
			},
		},
		"mappings": map[string]any{
			"properties": map[string]any{
				"Title": map[string]any{"type": "keyword"},
				"Headword": map[string]any{
					"type":     "text",
					"analyzer": "headword",
					"fields": map[string]any{
						"Smaller": map[string]any{
							"type":            "text",
							"analyzer":        "headword_smaller",
							"search_analyzer": "headword",
						},
					},
				},
				"HeadwordAlt": map[string]any{
					"type":     "text",
					"analyzer": "headword",
					"fields": map[string]any{
						"Smaller": map[string]any{
							"type":            "text",
							"analyzer":        "headword_smaller",
							"search_analyzer": "headword",
						},
					},
				},
				"Phrases": map[string]any{
					"type":     "text",
					"analyzer": "standard",
				},
				"Prefix": map[string]any{
					"type": "nested",
					"properties": map[string]any{
						"Letter1": map[string]any{"type": "keyword"},
						"Letter2": map[string]any{"type": "keyword"},
						"Letter3": map[string]any{"type": "keyword"},
						"Letter4": map[string]any{"type": "keyword"},
						"Letter5": map[string]any{"type": "keyword"},
					},
				},
				"Suggest": map[string]any{
					"type":                         "completion",
					"analyzer":                     "headword",
					"preserve_separators":          true,
					"preserve_position_increments": true,
					"max_input_length":             50,
				},
				"Content": map[string]any{
					"type":     "text",
					"analyzer": "body",
				},
				"ModifiedAt": map[string]any{
					"type": "date",
				},
			},
		},
	}, nil)

	if err != nil {
		return fmt.Errorf("storage put: %w", err)
	}
	return nil
}
