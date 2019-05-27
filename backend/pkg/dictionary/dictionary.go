package dictionary

// Dictionary represents a dictionary
type Dictionary interface {
	ID() string
	Title() string
	Slug() string
}

var dictionaries = []Dictionary{
	Rvblr{
		id:    "rvblr",
		title: "Тлумачальны слоўнік беларускай мовы (rv-blr.com)",
		slug:  "tlumachalny-slounik-bielaruskaj-movy-rv-blr-com",
	},
}
