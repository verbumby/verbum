package article

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/verbumby/verbum/backend/pkg/dictionary"
	"github.com/verbumby/verbum/backend/pkg/storage"
)

// Query queries articles in the storage
func Query(path string, reqbody interface{}) ([]Article, error) {
	respbody := struct {
		Hits struct {
			Hits []struct {
				Source Article `json:"_source"`
				Index  string  `json:"_index"`
			} `json:"hits"`
		} `json:"hits"`
	}{}

	if err := storage.Post("/dict-*/_search", reqbody, &respbody); err != nil {
		return nil, errors.Wrap(err, "query elastic")
	}

	result := []Article{}
	dicts := map[string]dictionary.Dictionary{}
	for _, hit := range respbody.Hits.Hits {
		dictID := strings.TrimPrefix(hit.Index, "dict-")
		if _, ok := dicts[dictID]; !ok {
			respbody := struct {
				Source dictionary.Dictionary `json:"_source"`
			}{}

			if err := storage.Get("/dicts/_doc/"+dictID, &respbody); err != nil {
				return nil, errors.Wrapf(err, "query dict %s", dictID)
			}

			dicts[dictID] = respbody.Source
		}

		article := hit.Source
		article.Dictionary = dicts[dictID]
		result = append(result, article)
	}

	return result, nil
}
