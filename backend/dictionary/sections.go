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
		Name: "Правераныя слоўнікі",
		DictIDs: []string{
			"tsblm2022", "tsbm", "klyshka", "sis1999", "sis2005", "bhn1971",
			"proverbia", "hsbm", "esbm", "brs", "rbs10", "abs", "susha", "pbs",
		},
	},
	{
		ID:      "dubious",
		Name:    "Слоўнік сумнеўнай якасці",
		DictIDs: []string{"beldeu", "kurjanka"},
		Descr:   "⚠️ Змешчаныя ў гэтым раздзеле слоўнікі не з'яўляюцца акадэмічнымі, не правераны намі і ўтрымліваюць крытычную колькасць памылак самага рознага роду (арфаграфічныя, у значэннях слоў і фраз, у націску). Карыстайцеся імі з асцярожнасцю!",
	},
}

func GetAllSections() []Section {
	return sections
}
