package textutil

var tableDeutsch = map[rune]any{
	'ä': "ae",
	'ö': "oe",
	'ü': "ue",
	'ß': "ss",
}

func SlugifyDeutsch(s string) string {
	return romanize(s, tableDeutsch)
}
