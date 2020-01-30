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
	Headwords    []string
	HeadwordsAlt []string
	Body         string
}
