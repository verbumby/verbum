package dictionary

import "html/template"

// Dictionary represents a dictionary
type Dictionary interface {
	ID() string
	Title() string
	ToHTML(content, title string) template.HTML
}

var dictionaries = []Dictionary{
	Markdown{
		id:    "rvblr",
		title: "Тлумачальны слоўнік беларускай літаратурнай мовы (2002, правапіс да 2008 г.)",
	},
	DSL{
		id:    "krapiva",
		title: "Тлумачальны слоўнік беларускай мовы (1977-84, правапіс да 2008 г.)",
	},
	DSL{
		id:    "brs",
		title: "Беларуска-рускі слоўнік, 4-е выданне (актуальны правапіс)",
	},
	DSL{
		id:                    "rbs",
		title:                 "Руска-беларускі слоўнік НАН РБ, 8-е выданне (правапіс да 2008 г.)",
		includeTitleInContent: true,
	},
	Stardict{
		id:    "rus-bel",
		title: "Руска-беларускі слоўнік НАН РБ, 8-е выданне (другая версія, правапіс да 2008 г.)",
	},
	DSL{
		id:                    "pashkievich",
		title:                 "Ангельска-беларускі слоўнік (В. Пашкевіч, 2006, класічны правапіс)",
		includeTitleInContent: true,
	},
	DSL{
		id:                    "hsbm",
		title:                 "Гістарычны слоўнік беларускай мовы (1982–2017)",
		includeTitleInContent: true,
	},
	DSL{
		id:    "esbm",
		title: "Этымалагічны слоўнік беларускай мовы, тамы 1-3 (1978-85, правапіс да 2008 г.)",
	},
}
