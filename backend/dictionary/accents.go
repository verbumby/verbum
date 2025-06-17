package dictionary

import "regexp"

var (
	reAccent           = regexp.MustCompile(`.\x{0301}`)
	accentSubstitution = `<span class="accent">$0</span>`
)

func wrapAccentedChars(s string) string {
	return reAccent.ReplaceAllString(s, accentSubstitution)
}
