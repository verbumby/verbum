package dslparser

import "fmt"

// DSLParser krapiva parser
func DSLParser(content string) (string, error) {
	result, err := Parse("article", []byte(content))
	if err != nil {
		return "", fmt.Errorf("parse article: %w", err)
	}
	return result.(string), nil
}
