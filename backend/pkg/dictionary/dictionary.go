package dictionary

// Dictionary represents a dictionary
type Dictionary struct {
	ID    string `json:"-"`
	Title string
	Slug  string
}
