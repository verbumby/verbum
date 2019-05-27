package article

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/verbumby/verbum/backend/pkg/dictionary"
	"github.com/verbumby/verbum/backend/pkg/storage"
)

// Query queries articles in the storage
func Query(path string, reqbody interface{}) ([]Article, int, error) {
	respbody := struct {
		Hits struct {
			Total int
			Hits  []struct {
				Source Article `json:"_source"`
				Index  string  `json:"_index"`
				ID     string  `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
	}{}

	if err := storage.Post(path, reqbody, &respbody); err != nil {
		return nil, 0, errors.Wrap(err, "query elastic")
	}

	result := []Article{}
	dicts := map[string]dictionary.Dictionary{}
	for _, hit := range respbody.Hits.Hits {
		dictID := strings.TrimPrefix(hit.Index, "dict-")
		if _, ok := dicts[dictID]; !ok {
			dict := dictionary.Get(dictID)
			if dict == nil {
				return nil, 0, fmt.Errorf("dictionary get %s: not found", dictID)
			}

			dicts[dictID] = dict
		}

		article := hit.Source
		article.ID = hit.ID
		article.Dictionary = dicts[dictID]
		result = append(result, article)
	}

	return result, respbody.Hits.Total, nil
}

// Get gets one article from the storage
func Get(dictionaryID, articleID string) (Article, error) {
	respbody := struct {
		Source Article `json:"_source"`
		Index  string  `json:"_index"`
		ID     string  `json:"_id"`
	}{}

	path := fmt.Sprintf("/dict-%s/_doc/%s", dictionaryID, articleID)
	if err := storage.Get(path, &respbody); err != nil {
		return Article{}, errors.Wrap(err, "storage get")
	}

	dict := dictionary.Get(dictionaryID)
	if dict == nil {
		return Article{}, fmt.Errorf("dictionary get %s: not found", dictionaryID)
	}

	article := respbody.Source
	article.Dictionary = dict
	article.ID = respbody.ID

	return article, nil
}
