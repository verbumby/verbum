package dictionary

import (
	"github.com/pkg/errors"
	"github.com/verbumby/verbum/backend/pkg/storage"
)

// Get gets a dictionary from storage
func Get(dictionaryID string) (Dictionary, error) {
	respbody := struct {
		ID     string     `json:"_id"`
		Source Dictionary `json:"_source"`
	}{}

	if err := storage.Get("/dicts/_doc/"+dictionaryID, &respbody); err != nil {
		return Dictionary{}, errors.Wrap(err, "storage get")
	}
	respbody.Source.ID = respbody.ID
	return respbody.Source, nil
}

// GetAll gets all dictionary from storage
func GetAll() ([]Dictionary, error) {
	respbody := struct {
		Hits struct {
			Hits []struct {
				ID     string     `json:"_id"`
				Source Dictionary `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}{}

	if err := storage.Post("/dicts/_doc/_search", map[string]interface{}{"size": 100}, &respbody); err != nil {
		return nil, errors.Wrap(err, "storage post")
	}

	result := []Dictionary{}
	for _, hit := range respbody.Hits.Hits {
		hit.Source.ID = hit.ID
		result = append(result, hit.Source)
	}

	return result, nil
}
