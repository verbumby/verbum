package utils

import (
	"strings"
	"unicode"
)

// Slugify returns slugged version of the string
func Slugify(s string) string {
	b := &strings.Builder{}

	s = Latinize(s)
	for _, r := range s {
		if unicode.IsUpper(r) {
			r = unicode.ToLower(r)
		}

		if unicode.Is(unicode.Latin, r) {
			b.WriteRune(r)
			continue
		}

		if unicode.IsDigit(r) {
			b.WriteRune(r)
			continue
		}
		if r == '(' || r == ')' || r == '-' || r == '.' || r == '\'' || r == ',' {
			b.WriteRune(r)
			continue
		}

		b.WriteRune('-')
	}
	return b.String()
}

var latinizeMap = map[rune]string{
	'І': "I",
	'і': "i",
	'Ў': "U",
	'ў': "u",
	'Ё': "YO",
	'ё': "yo",
	'А': "A",
	'а': "a",
	'Б': "B",
	'б': "b",
	'В': "V",
	'в': "v",
	'Г': "G",
	'г': "g",
	'Д': "D",
	'д': "d",
	'Е': "YE",
	'е': "ye",
	'Ж': "ZH",
	'ж': "zh",
	'З': "Z",
	'з': "z",
	'И': "I",
	'и': "i",
	'Й': "Y",
	'й': "y",
	'К': "K",
	'к': "k",
	'Л': "L",
	'л': "l",
	'М': "M",
	'м': "m",
	'Н': "N",
	'н': "n",
	'О': "O",
	'о': "o",
	'П': "P",
	'п': "p",
	'Р': "R",
	'р': "r",
	'С': "S",
	'с': "s",
	'Т': "T",
	'т': "t",
	'У': "U",
	'у': "u",
	'Ф': "F",
	'ф': "f",
	'Х': "KH",
	'х': "kh",
	'Ц': "TS",
	'ц': "ts",
	'Ч': "CH",
	'ч': "ch",
	'Ш': "SH",
	'ш': "sh",
	'Щ': "SHCH",
	'щ': "shch",
	'Ъ': "''",
	'ъ': "''",
	'Ы': "Y",
	'ы': "y",
	'Ь': "'",
	'ь': "'",
	'Э': "E",
	'э': "e",
	'Ю': "YU",
	'ю': "yu",
	'Я': "YA",
	'я': "ya",
}

// Latinize swaps non-latin runes with it's latin equivalents
func Latinize(s string) string {
	b := &strings.Builder{}
	for _, r := range s {
		if p, ok := latinizeMap[r]; ok {
			b.WriteString(p)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
