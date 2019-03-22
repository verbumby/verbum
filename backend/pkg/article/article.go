package article

import (
	"github.com/verbumby/verbum/backend/pkg/dictionary"
)

// Article represents an article entity
type Article struct {
	Dictionary dictionary.Dictionary
	Content    string
}
