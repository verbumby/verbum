package article

import (
	"errors"
	"fmt"

	"github.com/verbumby/verbum/backend/dictionary"
	"github.com/verbumby/verbum/backend/storage"
)

// Query queries articles in the storage
func Query(path string, reqbody interface{}) ([]Article, int, error) {
	respbody := struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source Article `json:"_source"`
				Index  string  `json:"_index"`
				ID     string  `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
	}{}

	if err := storage.Post(path, reqbody, &respbody); err != nil {
		return nil, 0, fmt.Errorf("query elastic: %w", err)
	}

	result := []Article{}
	for _, hit := range respbody.Hits.Hits {
		dict := dictionary.GetByIndex(hit.Index)
		if dict == nil {
			return nil, 0, fmt.Errorf("dictionary get %s: not found", hit.Index)
		}

		article := hit.Source
		article.ID = hit.ID
		article.Dictionary = dict
		result = append(result, article)
	}

	return result, respbody.Hits.Total.Value, nil
}

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
