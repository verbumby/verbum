package textutil

import (
	"regexp"
	"strings"
)

var belarusianConsonats = []rune{
	'Б', 'В', 'Г', 'Д', 'Ж', 'З', 'Й', 'К', 'Л', 'М', 'Н',
	'П', 'Р', 'С', 'Т', 'Ў', 'Ф', 'Х', 'Ц', 'Ч', 'Ш', 'Ь',
	'б', 'в', 'г', 'д', 'ж', 'з', 'й', 'к', 'л', 'м', 'н',
	'п', 'р', 'с', 'т', 'ў', 'ф', 'х', 'ц', 'ч', 'ш', 'ь',
}

func isBelarusianConsonat(r rune) bool {
	for _, c := range belarusianConsonats {
		if c == r {
			return true
		}
	}
	return false
}

func makeSoftBelarusianVowelRomanizer(hard, soft string) func(prev rune) string {
	return func(prev rune) string {
		if isBelarusianConsonat(prev) {
			return soft
		}

		return hard
	}
}

var table = map[rune]interface{}{
	'А': 'A',
	'а': 'a',
	'Б': 'B',
	'б': 'b',
	'В': 'V',
	'в': 'v',
	'Г': 'H',
	'г': 'h',
	'Д': 'D',
	'д': 'd',
	'Е': makeSoftBelarusianVowelRomanizer("Je", "ie"),
	'е': makeSoftBelarusianVowelRomanizer("je", "ie"),
	'Ё': makeSoftBelarusianVowelRomanizer("Jo", "io"),
	'ё': makeSoftBelarusianVowelRomanizer("Jo", "io"),
	'Ж': "Zh",
	'ж': "zh",
	'З': 'Z',
	'з': 'z',
	'І': 'I',
	'і': 'i',
	'Й': 'J',
	'й': 'j',
	'К': 'K',
	'к': 'k',
	'Л': 'L',
	'л': 'l',
	'М': 'M',
	'м': 'm',
	'Н': 'N',
	'н': 'n',
	'О': 'O',
	'о': 'o',
	'П': 'P',
	'п': 'p',
	'Р': 'R',
	'р': 'r',
	'С': 'S',
	'с': 's',
	'Т': 'T',
	'т': 't',
	'У': 'U',
	'у': 'u',
	'Ў': 'U',
	'ў': 'u',
	'Ф': 'F',
	'ф': 'f',
	'Х': "Ch",
	'х': "ch",
	'Ц': 'C',
	'ц': 'c',
	'Ч': "Ch",
	'ч': "ch",
	'Ш': "Sh",
	'ш': "sh",
	'’': nil,
	'Ы': 'Y',
	'ы': 'y',
	'Ь': nil,
	'ь': nil,
	'Э': 'E',
	'э': 'e',
	'Ю': makeSoftBelarusianVowelRomanizer("Ju", "iu"),
	'ю': makeSoftBelarusianVowelRomanizer("ju", "iu"),
	'Я': makeSoftBelarusianVowelRomanizer("Ja", "ia"),
	'я': makeSoftBelarusianVowelRomanizer("Ja", "ia"),
}

// RomanizeBelarusian romanizes belarusian text
func RomanizeBelarusian(s string) string {
	sb := &strings.Builder{}
	var prev rune
	for _, r := range s {
		v, ok := table[r]
		if !ok {
			sb.WriteRune(r)
			prev = r
			continue
		}

		switch v := v.(type) {
		case rune:
			sb.WriteRune(v)
		case string:
			sb.WriteString(v)
		case func(rune) string:
			sb.WriteString(v(prev))
		case nil:
			// do nothing;
		}
		prev = r
	}
	return sb.String()
}

// Slugify creates a url slug from the s string
func Slugify(s string) string {
	s = strings.ToLower(s)
	r := regexp.MustCompile("[^a-z0-9-]")
	s = r.ReplaceAllString(s, "-")

	for strings.Contains(s, "--") {
		s = strings.Replace(s, "--", "-", -1)
	}

	return s
}
