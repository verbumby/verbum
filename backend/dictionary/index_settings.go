package dictionary

type IndexSettings struct {
	PrependContentWithTitle          bool
	DictProvidesIDs                  bool
	DictProvidesIDsWithoutDuplicates bool
	LowercaseSuggestions             bool
}
