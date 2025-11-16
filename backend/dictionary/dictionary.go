package dictionary

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/verbumby/verbum/backend/config"
)

// Dictionary represents a dictionary
type Dictionary interface {
	ID() string
	IndexID() string
	Boost() float32
	Aliases() []string
	Title() string
	ToHTML(content string) template.HTML
	Preface() string
	Abbrevs() *Abbrevs
	ScanURL() string
	Slugifier() string
	Authors() bool
	IndexSettings() IndexSettings
}

var dictionaries []Dictionary

func InitDictionaries() error {
	abbrevs, err := loadDSLAbbrevs("grammardb/grammardb_abrv.dsl")
	if err != nil {
		return fmt.Errorf("load grammardb abbrevs: %w", err)
	}
	dictionaries = append(dictionaries, GrammarDB{
		Common: Common{
			id:        "grammardb",
			indexID:   "grammardb",
			boost:     1.3,
			title:     "Граматычная база Інстытута мовазнаўства НАН Беларусі (2025, актуальны правапіс)",
			abbrevs:   abbrevs,
			slugifier: "none",
			indexSettings: IndexSettings{
				DictProvidesIDs:                  true,
				DictProvidesIDsWithoutDuplicates: true,
			},
		},
	})

	abbrevs, err = loadDSLAbbrevs("tsblm2022/tsblm2022_abbr.txt")
	if err != nil {
		return fmt.Errorf("load tsblm2022 abbrevs: %w", err)
	}
	preface, err := loadPreface("tsblm2022/tsblm2022_pradmova.html")
	if err != nil {
		return fmt.Errorf("load tsblm2022 preface: %w", err)
	}
	dictionaries = append(dictionaries, HTML{
		Common: Common{
			id:        "tsblm2022",
			indexID:   "tsblm2022",
			boost:     1.6,
			title:     "Тлумачальны слоўнік беларускай літаратурнай мовы (І. Л. Капылоў, 2022, актуальны правапіс)",
			abbrevs:   abbrevs,
			preface:   preface,
			slugifier: "none",
			indexSettings: IndexSettings{
				DictProvidesIDs:                  true,
				DictProvidesIDsWithoutDuplicates: false,
			},
		},
	})

	abbrevs, err = loadDSLAbbrevs("tsblm/tsblm_abbr.txt")
	if err != nil {
		return fmt.Errorf("load tsblm abbrevs: %w", err)
	}
	dictionaries = append(dictionaries, HTML{
		Common: Common{
			id:        "tsblm",
			indexID:   "tsblm",
			boost:     1,
			aliases:   []string{"rvblr"},
			title:     "Тлумачальны слоўнік беларускай літаратурнай мовы (2002, правапіс да 2008 г.)",
			abbrevs:   abbrevs,
			slugifier: "belarusian",
			indexSettings: IndexSettings{
				DictProvidesIDs:                  true,
				DictProvidesIDsWithoutDuplicates: false,
			},
		},
	})

	abbrevs, err = loadDSLAbbrevs("tsbm/tsbm_abrv.dsl")
	if err != nil {
		return fmt.Errorf("load tsbm abbrevs: %w", err)
	}
	preface, err = loadPreface("tsbm/tsbm_pradmova.html")
	if err != nil {
		return fmt.Errorf("load tsbm preface: %w", err)
	}
	dictionaries = append(dictionaries, HTML{
		Common: Common{
			id:        "tsbm",
			indexID:   "tsbm",
			boost:     1.1,
			aliases:   []string{"krapiva"},
			title:     "Тлумачальны слоўнік беларускай мовы (1977-84, правапіс да 2008 г.)",
			preface:   preface,
			abbrevs:   abbrevs,
			scanURL:   "https://knihi.com/none/Tlumacalny_slounik_bielaruskaj_movy_zip.html",
			slugifier: "belarusian",
			indexSettings: IndexSettings{
				DictProvidesIDs:                  true,
				DictProvidesIDsWithoutDuplicates: true,
			},
		},
	})

	dictionaries = append(dictionaries, HTML{
		Common: Common{
			id:        "klyshka",
			indexID:   "klyshka",
			boost:     1,
			title:     "Слоўнік сінонімаў і блізказначных слоў, 2-е выданне (М. Клышка, правапіс да 2008 г.)",
			slugifier: "none",
			indexSettings: IndexSettings{
				DictProvidesIDs:                  true,
				DictProvidesIDsWithoutDuplicates: true,
			},
		},
	})

	abbrevs, err = loadDSLAbbrevs("sis1999/sis1999_abbr.txt")
	if err != nil {
		return fmt.Errorf("load sis1999 abbrevs: %w", err)
	}
	preface, err = loadPreface("sis1999/sis1999_pradmova.html")
	if err != nil {
		return fmt.Errorf("load sis1999 preface: %w", err)
	}
	dictionaries = append(dictionaries, HTML{
		Common: Common{
			id:        "sis1999",
			indexID:   "sis1999",
			boost:     1,
			title:     "Слоўнік іншамоўных слоў (А. Булыка, 1999, правапіс да 2008 г.)",
			preface:   preface,
			slugifier: "none",
			abbrevs:   abbrevs,
			scanURL:   "https://knihi.com/Alaksandar_Bulyka/Slounik_insamounych_slou_1999_pdf.zip.html",
			indexSettings: IndexSettings{
				DictProvidesIDs: true,
			},
		},
	})

	abbrevs, err = loadDSLAbbrevs("sis2005/sis2005_abbr.txt")
	if err != nil {
		return fmt.Errorf("load sis2005 abbrevs: %w", err)
	}
	dictionaries = append(dictionaries, HTML{
		Common: Common{
			id:        "sis2005",
			indexID:   "sis2005",
			boost:     1,
			title:     "Слоўнік іншамоўных слоў. Актуальная лексіка (А. Булыка, 2005, правапіс да 2008 г.)",
			slugifier: "none",
			abbrevs:   abbrevs,
			scanURL:   "https://knihi.com/Alaksandar_Bulyka/Slounik_insamounych_slou_Aktualnaja_leksika_2005.html",
			indexSettings: IndexSettings{
				DictProvidesIDs: true,
			},
		},
	})

	abbrevs, err = loadDSLAbbrevs("bhn1971/bhn1971_abbr.txt")
	if err != nil {
		return fmt.Errorf("load bhn1971 abbrevs: %w", err)
	}
	preface, err = loadPreface("bhn1971/bhn1971_pradmova.html")
	if err != nil {
		return fmt.Errorf("load bhn1971 preface: %w", err)
	}
	dictionaries = append(dictionaries, HTML{
		Common: Common{
			id:        "bhn1971",
			indexID:   "bhn1971",
			boost:     1,
			title:     "Беларускія геаграфічныя назвы. Тапаграфія. Гідралогія. (І. Яшкін, 1971, правапіс да 2008 г.)",
			slugifier: "none",
			preface:   preface,
			abbrevs:   abbrevs,
			scanURL:   "https://knihi.com/Ivan_Jaskin/Bielaruskija_hieahraficnyja_nazvy_1971.html",
			indexSettings: IndexSettings{
				DictProvidesIDs:                  true,
				DictProvidesIDsWithoutDuplicates: true,
			},
		},
	})

	preface, err = loadPreface("proverbia/proverbia_pradmova.html")
	if err != nil {
		return fmt.Errorf("load proverbia preface: %w", err)
	}
	dictionaries = append(dictionaries, HTML{
		Common: Common{
			id:        "proverbia",
			indexID:   "proverbia",
			boost:     1,
			title:     "Шасцімоўны слоўнік прыказак, прымавак і крылатых слоў (1993, правапіс да 2008 г.)",
			preface:   preface,
			slugifier: "none",
			indexSettings: IndexSettings{
				DictProvidesIDs:                  true,
				DictProvidesIDsWithoutDuplicates: true,
			},
		},
	})

	abbrevs, err = loadDSLAbbrevs("hsbm/hsbm_abrv.dsl")
	if err != nil {
		return fmt.Errorf("load hsbm abbrevs: %w", err)
	}
	preface, err = loadPreface("hsbm/hsbm_pradmova.html")
	if err != nil {
		return fmt.Errorf("load hsbm preface: %w", err)
	}
	dictionaries = append(dictionaries, DSL{
		Common: Common{
			id:        "hsbm",
			indexID:   "hsbm",
			boost:     0.9,
			title:     "Гістарычны слоўнік беларускай мовы (1982–2017)",
			preface:   preface,
			abbrevs:   abbrevs,
			scanURL:   "https://knihi.com/none/Histarycny_slounik_bielaruskaj_movy_zip.html",
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
			scanURL:   "https://knihi.com/none/Etymalahicny_slounik_bielaruskaj_movy_zip.html",
			slugifier: "none",
			indexSettings: IndexSettings{
				DictProvidesIDs:                  true,
				DictProvidesIDsWithoutDuplicates: true,
			},
		},
	})

	abbrevs, err = loadDSLAbbrevs("brs/brs_abrv.dsl")
	if err != nil {
		return fmt.Errorf("load brs abbrevs: %w", err)
	}
	preface, err = loadPreface("brs/brs_pradmova.html")
	if err != nil {
		return fmt.Errorf("load brs preface: %w", err)
	}
	dictionaries = append(dictionaries, DSL{
		Common: Common{
			id:        "brs",
			indexID:   "brs",
			boost:     1.3,
			title:     "Беларуска-рускі слоўнік, 4-е выданне (2012, актуальны правапіс)",
			preface:   preface,
			abbrevs:   abbrevs,
			slugifier: "belarusian",
		},
	})

	abbrevs, err = loadDSLAbbrevs("rbs10/rbs10_abbr.txt")
	if err != nil {
		return fmt.Errorf("load rbs10 abbrevs: %w", err)
	}
	preface, err = loadPreface("rbs10/rbs10_pradmova.html")
	if err != nil {
		return fmt.Errorf("load rbs10 preface: %w", err)
	}
	dictionaries = append(dictionaries, HTML{
		Common: Common{
			id:        "rbs10",
			indexID:   "rbs10",
			boost:     1.3,
			title:     "Руска-беларускі слоўнік НАН Беларусі, 10-е выданне (2012, актуальны правапіс)",
			preface:   preface,
			abbrevs:   abbrevs,
			slugifier: "none",
			indexSettings: IndexSettings{
				DictProvidesIDs:                  true,
				DictProvidesIDsWithoutDuplicates: true,
			},
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
		},
	})

	abbrevs, err = loadDSLAbbrevs("abs/abs_abrv.dsl")
	if err != nil {
		return fmt.Errorf("load abs abbrevs: %w", err)
	}
	preface, err = loadPreface("abs/abs_pradmova.html")
	if err != nil {
		return fmt.Errorf("load abs preface: %w", err)
	}
	dictionaries = append(dictionaries, DSL{
		Common: Common{
			id:      "abs",
			indexID: "abs",
			boost:   1,
			aliases: []string{"pashkievich"},
			title:   "Ангельска-беларускі слоўнік (В. Пашкевіч, 2006, класічны правапіс)",
			preface: preface,
			abbrevs: abbrevs,
			scanURL: "https://knihi.com/Valancina_Paskievic/Anhielska-bielaruski_slounik.html",
			authors: true,
			indexSettings: IndexSettings{
				PrependContentWithTitle: true,
			},
		},
	})

	abbrevs, err = loadDSLAbbrevs("susha/susha_abbr.txt")
	if err != nil {
		return fmt.Errorf("load susha abbrevs: %w", err)
	}
	preface, err = loadPreface("susha/susha_pradmova.html")
	if err != nil {
		return fmt.Errorf("load susha preface: %w", err)
	}
	dictionaries = append(dictionaries, HTML{
		Common: Common{
			id:      "susha",
			indexID: "susha",
			boost:   1.2,
			title:   "Англійска-беларускі слоўнік (Т. Суша, 2013, актуальны правапіс)",
			preface: preface,
			abbrevs: abbrevs,
			scanURL: "https://knihi.com/Tamara_Susa/Anhlijska-bielaruski_slounik.html",
			indexSettings: IndexSettings{
				DictProvidesIDs:                  true,
				DictProvidesIDsWithoutDuplicates: true,
			},
		},
	})

	abbrevs, err = loadDSLAbbrevs("pbs/pbs_abrv.dsl")
	if err != nil {
		return fmt.Errorf("load pbs abbrevs: %w", err)
	}
	preface, err = loadPreface("pbs/pbs_pradmova.html")
	if err != nil {
		return fmt.Errorf("load pbs preface: %w", err)
	}
	dictionaries = append(dictionaries, DSL{
		Common: Common{
			id:        "pbs",
			indexID:   "pbs",
			boost:     1,
			title:     "Польска-беларускі слоўнік (Я. Волкава, В. Авілава, 2004, правапіс да 2008 г.)",
			preface:   preface,
			slugifier: "polish",
			abbrevs:   abbrevs,
			scanURL:   "https://knihi.com/none/Polska-bielaruski_slounik_2004.html",
			indexSettings: IndexSettings{
				PrependContentWithTitle: true,
			},
		},
	})

	abbrevs, err = loadDSLAbbrevs("beldeu/beldeu_abbr.txt")
	if err != nil {
		return fmt.Errorf("load beldeu abbrevs: %w", err)
	}
	dictionaries = append(dictionaries, HTML{
		Common: Common{
			id:        "beldeu",
			indexID:   "beldeu",
			boost:     1,
			title:     "Беларуска-нямецкі слоўнік (М. Кур'янка, 2010, актуальны правапіс)",
			abbrevs:   abbrevs,
			scanURL:   "https://knihi.com/none/Bielaruska-niamiecki_slounik.html",
			slugifier: "none",
			authors:   true,
			indexSettings: IndexSettings{
				DictProvidesIDs:                  true,
				DictProvidesIDsWithoutDuplicates: true,
			},
		},
	})

	abbrevs, err = loadDSLAbbrevs("kurjanka/kurjanka_abrv.dsl")
	if err != nil {
		return fmt.Errorf("load kurjanka abbrevs: %w", err)
	}
	preface, err = loadPreface("kurjanka/kurjanka_pradmova.html")
	if err != nil {
		return fmt.Errorf("load kurjanka preface: %w", err)
	}
	dictionaries = append(dictionaries, HTML{
		Common: Common{
			id:      "kurjanka",
			indexID: "kurjanka",
			boost:   1,
			title:   "Нямецка-беларускі слоўнік (М. Кур'янка, 2006, правапіс да 2008 г.)",
			preface: preface,
			abbrevs: abbrevs,
			scanURL: "https://knihi.com/Mikalaj_Kurjanka/Niamiecka-bielaruski_slounik_2006.html",
			authors: true,
			indexSettings: IndexSettings{
				DictProvidesIDs:                  true,
				DictProvidesIDsWithoutDuplicates: true,
			},
		},
	})

	abbrevs, err = loadDSLAbbrevs("belen/bel_abbr.txt")
	if err != nil {
		return fmt.Errorf("load belen abbrevs: %w", err)
	}
	preface, err = loadPreface("belen/belen_pradmova.html")
	if err != nil {
		return fmt.Errorf("load belen preface: %w", err)
	}
	dictionaries = append(dictionaries, HTML{
		Common: Common{
			id:        "belen",
			indexID:   "belen",
			boost:     1,
			title:     "Беларуская Энцыклапедыя (1996—2004, правапіс да 2008 г., часткова)",
			preface:   preface,
			abbrevs:   abbrevs,
			slugifier: "none",
			scanURL:   "https://knihi.com/none/Bielaruskaja_encyklapiedyja_djvu.zip.html",
			indexSettings: IndexSettings{
				DictProvidesIDs:                  true,
				DictProvidesIDsWithoutDuplicates: false,
			},
		},
	})
	return nil
}

func loadPreface(path string) (string, error) {
	bytes, err := os.ReadFile(config.DictsRepoPath() + "/" + path)
	if err != nil {
		return "", fmt.Errorf("reading preface in %s failed: %w", path, err)
	}

	preface := string(bytes)

	if strings.Contains(preface, "</style>") {
		parts := strings.SplitN(preface, "</style>", 2)
		preface = parts[1]
	}

	return preface, nil
}
