package dictionary

import (
	"fmt"
	"html/template"
)

// Dictionary represents a dictionary
type Dictionary interface {
	ID() string
	IndexID() string
	Aliases() []string
	Title() string
	ToHTML(content string) template.HTML
	Abbrevs() *Abbrevs
}

var dictionaries []Dictionary

func InitDictionaries() error {
	dictionaries = append(dictionaries, Markdown{
		id:      "tsblm",
		indexID: "tsblm",
		aliases: []string{"rvblr"},
		title:   "Тлумачальны слоўнік беларускай літаратурнай мовы (2002, правапіс да 2008 г.)",
	})

	abbrevs, err := loadDSLAbbrevs("tsbm/tsbm_abrv.dsl")
	if err != nil {
		return fmt.Errorf("load tsbm abbrevs: %w", err)
	}
	dictionaries = append(dictionaries, DSL{
		id:      "tsbm",
		indexID: "tsbm",
		aliases: []string{"krapiva"},
		title:   "Тлумачальны слоўнік беларускай мовы (1977-84, правапіс да 2008 г.)",
		abbrevs: abbrevs,
	})

	abbrevs, err = loadDSLAbbrevs("hsbm/hsbm_abrv.dsl")
	if err != nil {
		return fmt.Errorf("load hsbm abbrevs: %w", err)
	}
	dictionaries = append(dictionaries, DSL{
		id:      "hsbm",
		indexID: "hsbm-2",
		title:   "Гістарычны слоўнік беларускай мовы (1982–2017, часткова)",
		abbrevs: abbrevs,
	})

	abbrevs, err = loadDSLAbbrevs("esbm/esbm_abrv.dsl")
	if err != nil {
		return fmt.Errorf("load esbm abbrevs: %w", err)
	}
	dictionaries = append(dictionaries, DSL{
		id:      "esbm",
		indexID: "esbm-6",
		title:   "Этымалагічны слоўнік беларускай мовы (1978-2017)",
		abbrevs: abbrevs,
	})

	abbrevs, err = loadDSLAbbrevs("brs/brs_abrv.dsl")
	if err != nil {
		return fmt.Errorf("load brs abbrevs: %w", err)
	}
	dictionaries = append(dictionaries, DSL{
		id:      "brs",
		indexID: "brs-2",
		title:   "Беларуска-рускі слоўнік, 4-е выданне (актуальны правапіс)",
		abbrevs: abbrevs,
	})

	abbrevs, err = loadDSLAbbrevs("rbs/rbs_abrv.dsl")
	if err != nil {
		return fmt.Errorf("load rbs abbrevs: %w", err)
	}
	dictionaries = append(dictionaries, DSL{
		id:      "rbs",
		indexID: "rbs-2",
		title:   "Руска-беларускі слоўнік НАН Беларусі, 8-е выданне (правапіс да 2008 г.)",
		abbrevs: abbrevs,
	})
	dictionaries = append(dictionaries, Stardict{
		id:      "rus-bel",
		indexID: "rus-bel-2",
		title:   "Руска-беларускі слоўнік НАН Беларусі, 8-е выданне (другая версія, правапіс да 2008 г.)",
	})
	dictionaries = append(dictionaries, DSL{
		id:      "abs",
		indexID: "abs-4",
		aliases: []string{"pashkievich"},
		title:   "Ангельска-беларускі слоўнік (В. Пашкевіч, 2006, класічны правапіс)",
	})
	dictionaries = append(dictionaries, HTML{
		id:      "susha",
		indexID: "susha-2",
		title:   "Англійска-беларускі слоўнік (Т. Суша, 2013, актуальны правапіс)",
	})
	dictionaries = append(dictionaries, DSL{
		id:      "pbs",
		indexID: "pbs",
		title:   "Польска-беларускі слоўнік (Я. Волкава, В. Авілава, 2004, правапіс да 2008 г.)",
	})
	dictionaries = append(dictionaries, DSL{
		id:      "kurjanka",
		indexID: "kurjanka-3",
		title:   "Нямецка-беларускі слоўнік (М. Кур'янка, 2006, правапіс да 2008 г.)",
	})
	return nil
}
