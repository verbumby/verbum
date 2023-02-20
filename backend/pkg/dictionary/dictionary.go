package dictionary

import "html/template"

// Dictionary represents a dictionary
type Dictionary interface {
	ID() string
	IndexID() string
	Aliases() []string
	Title() string
	ToHTML(content string) template.HTML
}

var dictionaries = []Dictionary{
	Markdown{
		id:      "tsblm",
		indexID: "tsblm",
		aliases: []string{"rvblr"},
		title:   "Тлумачальны слоўнік беларускай літаратурнай мовы (2002, правапіс да 2008 г.)",
	},
	DSL{
		id:      "tsbm",
		indexID: "tsbm",
		aliases: []string{"krapiva"},
		title:   "Тлумачальны слоўнік беларускай мовы (1977-84, правапіс да 2008 г.)",
	},
	DSL{
		id:      "hsbm",
		indexID: "hsbm",
		title:   "Гістарычны слоўнік беларускай мовы (1982–2017, часткова)",
	},
	DSL{
		id:      "esbm",
		indexID: "esbm-2",
		title:   "Этымалагічны слоўнік беларускай мовы (1978-2017)",
	},
	DSL{
		id:      "brs",
		indexID: "brs-2",
		title:   "Беларуска-рускі слоўнік, 4-е выданне (актуальны правапіс)",
	},
	DSL{
		id:      "rbs",
		indexID: "rbs",
		title:   "Руска-беларускі слоўнік НАН РБ, 8-е выданне (правапіс да 2008 г.)",
	},
	Stardict{
		id:      "rus-bel",
		indexID: "rus-bel-2",
		title:   "Руска-беларускі слоўнік НАН РБ, 8-е выданне (другая версія, правапіс да 2008 г.)",
	},
	DSL{
		id:      "abs",
		indexID: "abs-2",
		aliases: []string{"pashkievich"},
		title:   "Ангельска-беларускі слоўнік (В. Пашкевіч, 2006, класічны правапіс)",
	},
	HTML{
		id:      "susha",
		indexID: "susha-2",
		title:   "Англійска-беларускі слоўнік (Т. Суша, 2013, актуальны правапіс)",
	},
	DSL{
		id:      "pbs",
		indexID: "pbs",
		title:   "Польска-беларускі слоўнік (Я. Волкава, В. Авілава, 2004, правапіс да 2008 г.)",
	},
	DSL{
		id:      "kurjanka",
		indexID: "kurjanka-2",
		title:   "Нямецка-беларускі слоўнік (М. Кур'янка, 2006, правапіс да 2008 г.)",
	},
}
