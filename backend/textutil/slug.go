package textutil

import (
	"log"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var transformRemoveDiacritics = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

// Slugify creates a url slug from the s string
func Slugify(s string) string {
	var err error
	s, _, err = transform.String(transformRemoveDiacritics, s)
	if err != nil {
		log.Printf("failed to remove diacritics when slugifying %s", s)
	}

	r := regexp.MustCompile("[^a-zA-Z0-9-]")
	s = r.ReplaceAllString(s, "-")

	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}

	s = strings.Trim(s, "-")
	return s
}

func SlugifyLower(s string) string {
	return strings.ToLower(s)
}
