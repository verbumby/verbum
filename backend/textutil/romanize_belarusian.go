package textutil

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

var tableBelarusian = map[rune]interface{}{
	'А':  'A',
	'а':  'a',
	'Б':  'B',
	'б':  'b',
	'В':  'V',
	'в':  'v',
	'Г':  'H',
	'г':  'h',
	'Д':  'D',
	'д':  'd',
	'Е':  makeSoftBelarusianVowelRomanizer("Je", "ie"),
	'е':  makeSoftBelarusianVowelRomanizer("je", "ie"),
	'Ё':  makeSoftBelarusianVowelRomanizer("Jo", "io"),
	'ё':  makeSoftBelarusianVowelRomanizer("jo", "io"),
	'Ж':  "Zh",
	'ж':  "zh",
	'З':  'Z',
	'з':  'z',
	'І':  'I',
	'і':  'i',
	'Й':  'J',
	'й':  'j',
	'К':  'K',
	'к':  'k',
	'Л':  'L',
	'л':  'l',
	'М':  'M',
	'м':  'm',
	'Н':  'N',
	'н':  'n',
	'О':  'O',
	'о':  'o',
	'П':  'P',
	'п':  'p',
	'Р':  'R',
	'р':  'r',
	'С':  'S',
	'с':  's',
	'Т':  'T',
	'т':  't',
	'У':  'U',
	'у':  'u',
	'Ў':  'U',
	'ў':  'u',
	'Ф':  'F',
	'ф':  'f',
	'Х':  "Ch",
	'х':  "ch",
	'Ц':  'C',
	'ц':  'c',
	'Ч':  "Ch",
	'ч':  "ch",
	'Ш':  "Sh",
	'ш':  "sh",
	'’':  nil,
	'\'': nil,
	'Ы':  'Y',
	'ы':  'y',
	'Ь':  nil,
	'ь':  nil,
	'Э':  'E',
	'э':  'e',
	'Ю':  makeSoftBelarusianVowelRomanizer("Ju", "iu"),
	'ю':  makeSoftBelarusianVowelRomanizer("ju", "iu"),
	'Я':  makeSoftBelarusianVowelRomanizer("Ja", "ia"),
	'я':  makeSoftBelarusianVowelRomanizer("ja", "ia"),
}

var tableBelarusianV2 map[rune]any = map[rune]any{}

func init() {
	for k, v := range tableBelarusian {
		tableBelarusianV2[k] = v
	}

	tableBelarusianV2['Х'] = "Kh"
	tableBelarusianV2['х'] = "kh"
	tableBelarusianV2['Ь'] = "'"
	tableBelarusianV2['ь'] = "'"
}

func RomanizeBelarusian(s string) string {
	return romanize(s, tableBelarusian)
}

func RomanizeBelarusianV2(s string) string {
	return romanize(s, tableBelarusianV2)
}
