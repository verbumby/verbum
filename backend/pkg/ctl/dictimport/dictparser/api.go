package dictparser

type Dictionary struct {
	Meta        []MetaEntry
	Articles    []Article
	IDsProvided bool
}

type MetaEntry struct {
	Key   string
	Value string
}

type Article struct {
	ID           string
	Title        string
	Headwords    []string
	HeadwordsAlt []string
	Phrases      []string
	Body         string
}
