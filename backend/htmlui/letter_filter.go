package htmlui

import (
	"strconv"
)

// LetterFilter letter filter used on dictionary page
type LetterFilter struct {
	Prefix     []rune
	levels     [][]LetterFilterEntity
	LetterLink func(prefix string) string
}

// AddLevel adds filter level
func (lf *LetterFilter) AddLevel(level []LetterFilterEntity) {
	lf.levels = append(lf.levels, level)
}

// Links generates links
func (lf LetterFilter) Links() [][]LetterFilterLink {
	result := [][]LetterFilterLink{}
	for i, level := range lf.levels {
		result = append(result, []LetterFilterLink{})
		for _, e := range level {
			if e.Key == "" {
				continue
			}
			active := i < len(lf.Prefix) && string(lf.Prefix[i]) == e.Key
			prefix := string(lf.Prefix[:i])
			if !active {
				prefix += e.Key
			}
			result[i] = append(result[i], LetterFilterLink{
				URL:    lf.LetterLink(prefix),
				Text:   e.Key,
				Title:  strconv.FormatInt(int64(e.DocCount), 10),
				Active: active,
			})
		}
	}

	return result
}

// LetterFilterLink letter filter link
type LetterFilterLink struct {
	URL    string
	Text   string
	Active bool
	Title  string
}

// LetterFilterEntity letter filter entity
type LetterFilterEntity struct {
	Key      string `json:"key"`
	DocCount int    `json:"doc_count"`
}
