package dslparser

import (
	"github.com/pkg/errors"
)

// DSLParser krapiva parser
func DSLParser(content string) (string, error) {
	result, err := Parse("article", []byte(content))
	return result.(string), errors.Wrap(err, "parse article")
}
