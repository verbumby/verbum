package dictionary

import "html/template"

// Dictionary represents a dictionary
type Dictionary interface {
	ID() string
	Title() string
	Slug() string
	ToHTML(content, title string) template.HTML
}

var dictionaries = []Dictionary{
	Markdown{
		id:    "rvblr",
		title: "Тлумачальны слоўнік беларускай мовы (rv-blr.com)",
		slug:  "tlumachalny-slounik-bielaruskaj-movy-rv-blr-com",
	},
	DSL{
		id:    "krapiva",
		title: "Тлумачальны слоўнік беларускай мовы (Крапіва, 1977–1984)",
		slug:  "tlumachalny-slounik-bielaruskaj-movy-krapiva-1977–1984",
	},
	DSL{
		id:    "brs",
		title: "Беларуска-расійскі слоўнік, 4‑е выданне (Лукашанец, 2012)",
	},
	Stardict{
		id:    "bel-rus",
		title: "Беларуска-расійскі слоўнік",
		slug:  "bielaruska-rasijski-slounik",
	},
	DSL{
		id:                    "rbs",
		title:                 "Расійска-беларускі слоўнік, 8-е выданне (Крапіва, 2002)",
		includeTitleInContent: true,
	},
	Stardict{
		id:    "rus-bel",
		title: "Расійска-беларускі універсальны слоўнік",
		slug:  "rasijska-bielaruski-universalny-slounik",
	},
	DSL{
		id:                    "pashkievich",
		title:                 "Ангельска-беларускі слоўнік (Пашкевіч, 2006)",
		slug:                  "anhyelska-byelaruski-slounik-pashkyevich-2006",
		includeTitleInContent: true,
	},
	DSL{
		id:                    "hsbm",
		title:                 "Гістарычны слоўнік беларускай мовы (1982–2017)",
		slug:                  "histarychny-slounik-bielaruskaj-movy-1982-2017",
		includeTitleInContent: true,
	},
	DSL{
		id:    "esbm",
		title: "Этымалагічны слоўнік беларускай мовы",
	},
}
