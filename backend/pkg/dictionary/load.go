package dictionary

import (
	"github.com/pkg/errors"
	"github.com/verbumby/verbum/backend/pkg/storage"
)

// Get gets a dictionary from storage
func Get(dictionaryID string) (Dictionary, error) {
	respbody := struct {
		Source Dictionary `json:"_source"`
	}{}

	if err := storage.Get("/dicts/_doc/"+dictionaryID, &respbody); err != nil {
		return Dictionary{}, errors.Wrap(err, "storage get")
	}
	return respbody.Source, nil
}
