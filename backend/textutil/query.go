package textutil

import "strings"

func NormalizeQuery(q string) string {
	if strings.ContainsRune(q, '\'') {
		q = strings.ReplaceAll(q, "'", "’")
	}
	return q
}
