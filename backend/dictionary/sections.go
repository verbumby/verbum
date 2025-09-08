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
		Name: "Рэкамендаваныя слоўнікі",
		DictIDs: []string{
			"tsblm2022", "tsbm", "klyshka", "sis1999", "sis2005", "bhn1971",
			"proverbia", "hsbm", "esbm", "brs", "rbs10", "abs", "susha", "pbs",
		},
	},
	{
		ID:      "authors",
		Name:    "Аўтарскія слоўнікі",
		DictIDs: []string{"beldeu", "kurjanka"},
		Descr:   "⚠️ Змешчаныя ў гэтым раздзеле слоўнікі з'яўляюцца аўтарскімі, не правераны намі і ўтрымліваюць крытычную колькасць памылак самага рознага роду (арфаграфічныя, у значэннях слоў і фраз, у націску). Карыстайцеся імі з асцярожнасцю!",
	},
}

func GetAllSections() []Section {
	return sections
}
