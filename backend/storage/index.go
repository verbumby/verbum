package storage

import "fmt"

func CreateDictIndex(dictID string) error {
	err := Put("/dict-"+dictID, map[string]any{
		"settings": map[string]any{
			"number_of_shards":   1,
			"number_of_replicas": 0,
			"max_result_window":  400_000,
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
						"filter":      []string{"lowercase", "strip_diacritics"},
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
				"filter": map[string]any{
					"strip_diacritics": map[string]any{
						"type": "icu_transform",
						"id":   "NFD; [\\u0301\\u030C\\u0311] Remove; NFC;",
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
				"Content": map[string]any{
					"type":     "text",
					"analyzer": "body",
				},
			},
		},
	}, nil)
	if err != nil {
		return fmt.Errorf("storage put: %w", err)
	}
	return nil
}

func CreateSuggestIndex(dictID string) error {
	err := Put("/sugg-"+dictID, map[string]any{
		"settings": map[string]any{
			"number_of_shards":   1,
			"number_of_replicas": 0,
			"analysis": map[string]any{
				"analyzer": map[string]any{
					"headword": map[string]any{
						"filter":    []string{"lowercase"},
						"type":      "custom",
						"tokenizer": "keyword",
					},
				},
			},
		},
		"mappings": map[string]any{
			"properties": map[string]any{
				"Suggest": map[string]any{
					"type":                         "completion",
					"analyzer":                     "headword",
					"preserve_separators":          true,
					"preserve_position_increments": true,
					"max_input_length":             50,
				},
			},
		},
	}, nil)
	if err != nil {
		return fmt.Errorf("storage put: %w", err)
	}
	return nil
}
