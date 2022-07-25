package dictionary

import "golang.org/x/exp/slices"

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

func GetByIndexID(indexID string) Dictionary {
	for _, d := range dictionaries {
		if indexID == d.IndexID() {
			return d
		}
	}
	return nil
}
