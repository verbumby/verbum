package article

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

func parseArticle(content string) ([][]html.Token, error) {
	result := [][]html.Token{}
	queue := []html.Token{}
	z := html.NewTokenizer(strings.NewReader(content))
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			err := z.Err()
			if err != io.EOF {
				return nil, errors.Wrap(err, "tokenize html")
			}
			break
		}

		t := z.Token()
		queue = append(queue, t)
		if tt == html.EndTagToken {
			startingTagIDX := len(queue) - 1
			for startingTagIDX >= 0 && queue[startingTagIDX].Type != html.StartTagToken {
				startingTagIDX--
			}
			if startingTagIDX < 0 {
				return nil, fmt.Errorf("no starting tag found for closing tag %s", t.Type)
			}
			if err := validateTagTokens(queue[startingTagIDX:]); err != nil {
				return nil, errors.Wrap(err, "processing tags sequence")
			}

			result = append(result, queue[startingTagIDX:])
		}
	}

	return result, nil
}

func validateTagTokens(tokens []html.Token) error {
	if len(tokens) < 2 {
		return fmt.Errorf("expected at least two (openning and closing) tokens")
	}

	qs := tokens[0]
	qe := tokens[len(tokens)-1]

	if qs.Type != html.StartTagToken {
		return fmt.Errorf("expected starting token type to be %s, got %s", html.StartTagToken, qs.Type)
	}

	if tokens[len(tokens)-1].Type != html.EndTagToken {
		return fmt.Errorf("expected ending token type to be %s, got %s", html.EndTagToken, qe.Type)
	}

	// check whether start and end tokens match
	if qs.Data != qe.Data {
		return fmt.Errorf("closing tag %s doesn't match opening tag %s", qe.Data, qs.Data)
	}

	return nil
}

// RvblrParse returns rvblr article headwords struct
func RvblrParse(content string) (hws []string, hwsalt []string, err error) {
	parse := func(content string) ([]string, error) {
		tokenGroups, err := parseArticle(content)
		if err != nil {
			return nil, errors.Wrap(err, "get html tokens")
		}

		headwordTokenGroups := [][]html.Token{}
		for _, tokenGroup := range tokenGroups {
			if tokenGroup[0].Data == "v-hw" {
				headwordTokenGroups = append(headwordTokenGroups, tokenGroup)
			}
		}

		result := []string{}
		for _, tokens := range headwordTokenGroups {
			hw := ""
			for _, t := range tokens[1 : len(tokens)-1] {
				hw += t.Data
			}
			result = append(result, hw)
		}

		return result, nil
	}
	contents := strings.SplitN(content, "\n", 2)
	fmt.Println(contents[1])

	hws, err = parse(contents[0])
	if err != nil {
		err = errors.Wrap(err, "parse headwords")
	}

	hwsalt, err = parse(contents[1])
	if err != nil {
		err = errors.Wrap(err, "parse alt headwords")
	}

	return
}
