package dictparser

type Dictionary struct {
	Meta     []MetaEntry
	Articles []Article
}

type MetaEntry struct {
	Key   string
	Value string
}

type Article struct {
	ID           string
	Headwords    []string
	HeadwordsAlt []string
	Phrases      []string
	Body         string
}
