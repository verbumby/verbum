package dictimport

import (
	"fmt"
	"log"
	"strings"

	"github.com/verbumby/verbum/backend/dictionary"
	"github.com/verbumby/verbum/backend/textutil"
)

type IDGen struct {
	dict  dictionary.Dictionary
	cache map[string]int
	dups  map[string]int
}

func NewIDGen(dict dictionary.Dictionary) *IDGen {
	return &IDGen{
		dict:  dict,
		cache: map[string]int{},
		dups:  map[string]int{},
	}
}

func (ig *IDGen) Gen(firstHW, articleID string) (string, error) {
	id := strings.ToLower(firstHW)
	if ig.dict.IndexSettings().DictProvidesIDs {
		id = articleID
		ig.dups[id]++
		if ig.dict.IndexSettings().DictProvidesIDsWithoutDuplicates && ig.dups[id] > 1 {
			return "", fmt.Errorf("duplicate id: %s", id)
		}
	}
	var err error
	id, err = ig.assembleID(id)
	if err != nil {
		return "", err
	}

	idbase, idbaseno := calcIDBase(id)

	ig.cache[id]++
	if ig.cache[id] > 1 || (idbaseno > -1 && ig.cache[idbase] >= idbaseno) {
		id = fmt.Sprintf("%s-%d", id, ig.cache[id])
		log.Printf("adding index to id %s", id)
	}
	return id, nil
}

func (ig *IDGen) assembleID(firstHW string) (string, error) {
	hw := firstHW
	var romanized string
	switch ig.dict.Slugifier() {
	case "belarusian":
		romanized = textutil.RomanizeBelarusian(hw)
	case "none":
		return firstHW, nil
	case "russian":
		romanized = textutil.RomanizeRussian(hw)
	case "polish":
		romanized = textutil.SlugifyPolish(hw)
	case "german":
		romanized = textutil.SlugifyDeutsch(hw)
	case "":
		romanized = hw
	default:
		return "", fmt.Errorf("unknown romanizing strategy: %s", ig.dict.Slugifier())
	}
	result := romanized
	return textutil.Slugify(result), nil
}
