package dictionary

import (
	"html/template"
	"strings"
)

// Stardict dictionary
type Stardict struct {
	Common
}

// ToHTML implements Dictionary interface
func (d Stardict) ToHTML(content string) template.HTML {
	content = strings.ReplaceAll(content, "<k>", "<p><v-hw>")
	content = strings.ReplaceAll(content, "</k>", "</v-hw></p>")
	return template.HTML(content)
}
