package dictionary

type Section struct {
	ID      string
	Name    string
	DictIDs []string
	Descr   string
}

var sections = []Section{
	{
		ID:   "default",
		Name: "Асноўныя слоўнікі",
		DictIDs: []string{
			"tsblm2022", "tsbm", "klyshka", "sis1999", "sis2005", "bhn1971",
			"proverbia", "hsbm", "esbm", "brs", "rbs10", "abs", "susha", "pbs",
		},
	},
	{
		ID:      "deu",
		Name:    "Слоўнікі нямецкай мовы",
		DictIDs: []string{"beldeu", "kurjanka"},
		Descr:   "⚠️ Гэтыя слоўнікі не з'яўляюцца акадэмічнымі і ўтрымліваюць мноства памылак.",
	},
}

func GetAllSections() []Section {
	return sections
}
