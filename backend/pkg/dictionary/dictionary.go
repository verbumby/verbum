package dictionary

import "html/template"

// Dictionary represents a dictionary
type Dictionary interface {
	ID() string
	Title() string
	Slug() string
	ToHTML(content string) template.HTML
}

var dictionaries = []Dictionary{
	Rvblr{
		id:    "rvblr",
		title: "Тлумачальны слоўнік беларускай мовы (rv-blr.com)",
		slug:  "tlumachalny-slounik-bielaruskaj-movy-rv-blr-com",
	},
	Krapiva{
		id:    "krapiva",
		title: "Тлумачальны слоўнік беларускай мовы (Крапіва, 1977–1984)",
		slug:  "tlumachalny-slounik-bielaruskaj-movy-krapiva-1977–1984",
	},
}
