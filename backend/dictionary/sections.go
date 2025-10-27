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
		Name: "Усе",
		DictIDs: []string{
			"grammardb", "tsblm2022", "tsbm", "klyshka", "sis1999", "sis2005", "bhn1971",
			"proverbia", "hsbm", "esbm", "belen", "brs", "rbs10", "abs", "susha", "pbs",
			"beldeu", "kurjanka",
		},
	},
	{
		ID:      "main",
		Name:    "Асноўныя",
		DictIDs: []string{"grammardb", "brs", "rbs10", "tsbm", "tsblm2022", "esbm", "klyshka"},
		Descr:   "Самы неабходны мінімум слоўнікаў беларускай мовы. Нічога лішняга.",
	},
	{
		ID:      "authors",
		Name:    "Аўтарскія",
		DictIDs: []string{"abs", "beldeu", "kurjanka"},
		Descr:   "⚠️ Тут прадстаўлены аўтарскія слоўнікі — у іх словы і тлумачэнні пададзены паводле асабістых поглядаў укладальнікаў. Магчымыя няправільныя націскі, а таксама іншыя памылкі і недакладнасці.",
	},
	{
		ID:      "encyclopedias",
		Name:    "Энцыклапедыі",
		DictIDs: []string{"belen"},
		Descr:   "Тут знаходзяцца артыкулы з «Беларускай Энцыклапедыі» — афіцыйнай крыніцы дакладнай і навукова праверанай інфармацыі.",
	},
}

func GetAllSections() []Section {
	return sections
}
