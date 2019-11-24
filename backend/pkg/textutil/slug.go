package textutil

import (
	"regexp"
	"strings"
)

// Slugify creates a url slug from the s string
func Slugify(s string) string {
	s = strings.ToLower(s)
	r := regexp.MustCompile("[^a-z0-9-]")
	s = r.ReplaceAllString(s, "-")

	for strings.Contains(s, "--") {
		s = strings.Replace(s, "--", "-", -1)
	}

	s = strings.Trim(s, "-")
	return s
}
