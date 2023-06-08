package textutil

var tablePolish = map[rune]any{
	'ą': 'a',
	'ć': 'c',
	'ę': 'e',
	'ł': 'l',
	'ń': 'n',
	'ó': 'o',
	'ś': 's',
	'ź': 'z',
	'ż': 'z',
}

func SlugifyPolish(s string) string {
	return romanize(s, tablePolish)
}
