package textutil

import "strings"

func romanize(s string, table map[rune]interface{}) string {
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
