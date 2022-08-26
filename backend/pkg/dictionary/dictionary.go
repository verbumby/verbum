package dictionary

import "html/template"

// Dictionary represents a dictionary
type Dictionary interface {
	ID() string
	IndexID() string
	Aliases() []string
	Title() string
	ToHTML(content, title string) template.HTML
}

var dictionaries = []Dictionary{
	Markdown{
		id:      "tsblm",
		indexID: "rvblr",
		aliases: []string{"rvblr"},
		title:   "Тлумачальны слоўнік беларускай літаратурнай мовы (2002, правапіс да 2008 г.)",
	},
	DSL{
		id:      "tsbm",
		indexID: "krapiva",
		aliases: []string{"krapiva"},
		title:   "Тлумачальны слоўнік беларускай мовы (1977-84, правапіс да 2008 г.)",
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
		id:                    "abs",
		indexID:               "pashkievich",
		aliases:               []string{"pashkievich"},
		title:                 "Ангельска-беларускі слоўнік (В. Пашкевіч, 2006, класічны правапіс)",
		includeTitleInContent: true,
	},
	HTML{
		id:      "susha",
		indexID: "susha-2",
		title:   "Англійска-беларускі слоўнік (Т. Суша, 2013)",
	},
	DSL{
		id:                    "hsbm",
		title:                 "Гістарычны слоўнік беларускай мовы (1982–2017, часткова)",
		includeTitleInContent: true,
	},
	DSL{
		id:    "esbm",
		title: "Этымалагічны слоўнік беларускай мовы (1978-2017)",
	},
}
