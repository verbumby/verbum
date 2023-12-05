package dictionary

import (
	"regexp"
	"slices"
)

func Get(idOrAlias string) Dictionary {
	var result Dictionary
	for _, d := range dictionaries {
		if d.ID() == idOrAlias || slices.Contains(d.Aliases(), idOrAlias) {
			result = d
			break
		}
	}

	return result
}

func GetAll() []Dictionary {
	return dictionaries
}

func GetAllListed() []Dictionary {
	result := []Dictionary{}
	for _, d := range GetAll() {
		if !d.Unlisted() {
			result = append(result, d)
		}
	}
	return result
}

func GetByID(id string) Dictionary {
	for _, d := range dictionaries {
		if id == d.ID() {
			return d
		}
	}
	return nil
}

var reIndexToID = regexp.MustCompile(`^(?:dict|sugg)-(.+?)(!?-\d+)?$`)

func GetByIndex(index string) Dictionary {
	match := reIndexToID.FindStringSubmatch(index)
	return GetByID(match[1])
}
