package dictionary

import (
	"fmt"
	"html/template"
)

// Dictionary represents a dictionary
type Dictionary interface {
	ID() string
	IndexID() string
	Boost() float32
	Aliases() []string
	Title() string
	ToHTML(content string) template.HTML
	Abbrevs() *Abbrevs
	Slugifier() string
	Unlisted() bool
	IndexSettings() IndexSettings
}

var dictionaries []Dictionary

func InitDictionaries() error {
	dictionaries = append(dictionaries, HTML{
		Common: Common{
			id:        "tsblm",
			indexID:   "tsblm",
			boost:     1.6,
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
			boost:     1.1,
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
			id:        "hsbm",
			indexID:   "hsbm",
			boost:     0.9,
			title:     "Гістарычны слоўнік беларускай мовы (1982–2017, часткова)",
			abbrevs:   abbrevs,
			slugifier: "russian",
			indexSettings: IndexSettings{
				PrependContentWithTitle: true,
			},
		},
	})

	abbrevs, err = loadDSLAbbrevs("esbm/esbm_abrv.dsl")
	if err != nil {
		return fmt.Errorf("load esbm abbrevs: %w", err)
	}
	dictionaries = append(dictionaries, HTML{
		Common: Common{
			id:        "esbm",
			indexID:   "esbm",
			boost:     0.9,
			title:     "Этымалагічны слоўнік беларускай мовы (1978-2017)",
			abbrevs:   abbrevs,
			slugifier: "none",
			indexSettings: IndexSettings{
				ConvertHeadwordsToLowercase: true,
			},
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
			boost:     1.3,
			title:     "Беларуска-рускі слоўнік, 4-е выданне (актуальны правапіс)",
			abbrevs:   abbrevs,
			slugifier: "belarusian",
		},
	})

	abbrevs, err = loadDSLAbbrevs("rbs10/rbs10_abbr.txt")
	if err != nil {
		return fmt.Errorf("load rbs10 abbrevs: %w", err)
	}
	dictionaries = append(dictionaries, HTML{
		Common: Common{
			id:        "rbs10",
			indexID:   "rbs10",
			boost:     1.3,
			title:     "Руска-беларускі слоўнік НАН Беларусі, 10-е выданне (актуальны правапіс)",
			abbrevs:   abbrevs,
			slugifier: "none",
		},
	})

	abbrevs, err = loadDSLAbbrevs("rbs/rbs_abrv.dsl")
	if err != nil {
		return fmt.Errorf("load rbs abbrevs: %w", err)
	}
	dictionaries = append(dictionaries, DSL{
		Common: Common{
			id:        "rbs",
			indexID:   "rbs",
			boost:     1,
			title:     "Руска-беларускі слоўнік НАН Беларусі, 8-е выданне (правапіс да 2008 г.)",
			abbrevs:   abbrevs,
			slugifier: "russian",
			unlisted:  true,
			indexSettings: IndexSettings{
				PrependContentWithTitle: true,
			},
		},
	})
	dictionaries = append(dictionaries, Stardict{
		Common: Common{
			id:        "rus-bel",
			indexID:   "rus-bel",
			boost:     1,
			title:     "Руска-беларускі слоўнік НАН Беларусі, 8-е выданне (другая версія, правапіс да 2008 г.)",
			slugifier: "russian",
			unlisted:  true,
		},
	})
	dictionaries = append(dictionaries, DSL{
		Common: Common{
			id:      "abs",
			indexID: "abs",
			boost:   1,
			aliases: []string{"pashkievich"},
			title:   "Ангельска-беларускі слоўнік (В. Пашкевіч, 2006, класічны правапіс)",
			indexSettings: IndexSettings{
				PrependContentWithTitle: true,
			},
		},
	})
	dictionaries = append(dictionaries, HTML{
		Common: Common{
			id:      "susha",
			indexID: "susha",
			boost:   1,
			title:   "Англійска-беларускі слоўнік (Т. Суша, 2013, актуальны правапіс)",
		},
	})
	dictionaries = append(dictionaries, DSL{
		Common: Common{
			id:        "pbs",
			indexID:   "pbs",
			boost:     1,
			title:     "Польска-беларускі слоўнік (Я. Волкава, В. Авілава, 2004, правапіс да 2008 г.)",
			slugifier: "polish",
			indexSettings: IndexSettings{
				PrependContentWithTitle: true,
			},
		},
	})
	dictionaries = append(dictionaries, DSL{
		Common: Common{
			id:        "kurjanka",
			indexID:   "kurjanka",
			boost:     1,
			title:     "Нямецка-беларускі слоўнік (М. Кур'янка, 2006, правапіс да 2008 г.)",
			slugifier: "german",
			indexSettings: IndexSettings{
				PrependContentWithTitle: true,
			},
		},
	})
	return nil
}
