package krapivaparser

import (
	"github.com/pkg/errors"
)

// KrapivaParser krapiva parser
func KrapivaParser(content string) (string, error) {
	result, err := Parse("article", []byte(content))
	return result.(string), errors.Wrap(err, "parse article")
}
