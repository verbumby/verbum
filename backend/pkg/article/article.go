package article

import (
	"github.com/verbumby/verbum/backend/pkg/dictionary"
)

// Article represents an article entity
type Article struct {
	ID          string                `json:"-"`
	Dictionary  dictionary.Dictionary `json:"-"`
	Title       string
	Content     string
	Headword    []string
	HeadwordAlt []string
	Suggest     []Suggest
	Prefix      []Prefix
}

// Suggest article suggests
type Suggest struct {
	Input  string `json:"input"`
	Weight int    `json:"weight"`
}

// Prefix article prefix struct
type Prefix struct {
	Letter1 string
	Letter2 string
	Letter3 string
}
