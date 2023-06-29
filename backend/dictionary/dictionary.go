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
	PrependContentWithTitle() bool
	Slugifier() string
}

var dictionaries []Dictionary

func InitDictionaries() error {
	dictionaries = append(dictionaries, HTML{
		Common: Common{
			id:        "tsblm",
			indexID:   "tsblm",
			aliases:   []string{"rvblr"},
			title:     "Тлумачальны слоўнік беларускай літаратурнай мовы (2002, правапіс да 2008 г.)",
			slugifier: "belarusian",
		},
	})

	abbrevs, err := loadDSLAbbrevs("tsbm/tsbm_abrv.dsl")
	if err != nil {
		return fmt.Errorf("load tsbm abbrevs: %w", err)
	}
	dictionaries = append(dictionaries, DSL{
		Common: Common{
			id:        "tsbm",
			indexID:   "tsbm",
			aliases:   []string{"krapiva"},
			title:     "Тлумачальны слоўнік беларускай мовы (1977-84, правапіс да 2008 г.)",
			abbrevs:   abbrevs,
			slugifier: "belarusian",
		},
	})

	abbrevs, err = loadDSLAbbrevs("hsbm/hsbm_abrv.dsl")
	if err != nil {
		return fmt.Errorf("load hsbm abbrevs: %w", err)
	}
	dictionaries = append(dictionaries, DSL{
		Common: Common{
			id:                      "hsbm",
			indexID:                 "hsbm",
			title:                   "Гістарычны слоўнік беларускай мовы (1982–2017, часткова)",
			abbrevs:                 abbrevs,
			prependContentWithTitle: true,
			slugifier:               "russian",
		},
	})

	abbrevs, err = loadDSLAbbrevs("esbm/esbm_abrv.dsl")
	if err != nil {
		return fmt.Errorf("load esbm abbrevs: %w", err)
	}
	dictionaries = append(dictionaries, DSL{
		Common: Common{
			id:        "esbm",
			indexID:   "esbm",
			title:     "Этымалагічны слоўнік беларускай мовы (1978-2017)",
			abbrevs:   abbrevs,
			slugifier: "belarusian",
		},
	})

	abbrevs, err = loadDSLAbbrevs("brs/brs_abrv.dsl")
	if err != nil {
		return fmt.Errorf("load brs abbrevs: %w", err)
	}
	dictionaries = append(dictionaries, DSL{
		Common: Common{
			id:        "brs",
			indexID:   "brs",
			title:     "Беларуска-рускі слоўнік, 4-е выданне (актуальны правапіс)",
			abbrevs:   abbrevs,
			slugifier: "belarusian",
		},
	})

	abbrevs, err = loadDSLAbbrevs("rbs/rbs_abrv.dsl")
	if err != nil {
		return fmt.Errorf("load rbs abbrevs: %w", err)
	}
	dictionaries = append(dictionaries, DSL{
		Common: Common{
			id:                      "rbs",
			indexID:                 "rbs",
			title:                   "Руска-беларускі слоўнік НАН Беларусі, 8-е выданне (правапіс да 2008 г.)",
			abbrevs:                 abbrevs,
			prependContentWithTitle: true,
			slugifier:               "russian",
		},
	})
	dictionaries = append(dictionaries, Stardict{
		Common: Common{
			id:        "rus-bel",
			indexID:   "rus-bel",
			title:     "Руска-беларускі слоўнік НАН Беларусі, 8-е выданне (другая версія, правапіс да 2008 г.)",
			slugifier: "russian",
		},
	})
	dictionaries = append(dictionaries, DSL{
		Common: Common{
			id:                      "abs",
			indexID:                 "abs",
			aliases:                 []string{"pashkievich"},
			title:                   "Ангельска-беларускі слоўнік (В. Пашкевіч, 2006, класічны правапіс)",
			prependContentWithTitle: true,
		},
	})
	dictionaries = append(dictionaries, HTML{
		Common: Common{
			id:      "susha",
			indexID: "susha",
			title:   "Англійска-беларускі слоўнік (Т. Суша, 2013, актуальны правапіс)",
		},
	})
	dictionaries = append(dictionaries, DSL{
		Common: Common{
			id:                      "pbs",
			indexID:                 "pbs",
			title:                   "Польска-беларускі слоўнік (Я. Волкава, В. Авілава, 2004, правапіс да 2008 г.)",
			slugifier:               "polish",
			prependContentWithTitle: true,
		},
	})
	dictionaries = append(dictionaries, DSL{
		Common: Common{
			id:                      "kurjanka",
			indexID:                 "kurjanka",
			title:                   "Нямецка-беларускі слоўнік (М. Кур'янка, 2006, правапіс да 2008 г.)",
			prependContentWithTitle: true,
			slugifier:               "german",
		},
	})
	return nil
}
