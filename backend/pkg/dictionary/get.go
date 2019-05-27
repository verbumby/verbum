package dictionary

// Get a dictionary by ID
func Get(dictionaryID string) Dictionary {
	var result Dictionary
	for _, d := range dictionaries {
		if d.ID() == dictionaryID {
			result = d
			break
		}
	}

	return result
}

// GetAll dictionaries
func GetAll() []Dictionary {
	return dictionaries
}
