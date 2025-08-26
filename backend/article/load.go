package article

import (
	"errors"
	"fmt"

	"github.com/verbumby/verbum/backend/dictionary"
	"github.com/verbumby/verbum/backend/storage"
)

// Get gets one article from the storage
func Get(d dictionary.Dictionary, articleID string) (Article, error) {
	respbody := struct {
		Source Article `json:"_source"`
		Index  string  `json:"_index"`
		ID     string  `json:"_id"`
	}{}

	path := fmt.Sprintf("/dict-%s/_doc/%s", d.IndexID(), articleID)
	if err := storage.Get(path, &respbody); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return Article{}, nil
		}
		return Article{}, fmt.Errorf("storage get: %w", err)
	}

	article := respbody.Source
	article.Dictionary = d
	article.ID = respbody.ID

	return article, nil
}
