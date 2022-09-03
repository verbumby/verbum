package html

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/verbumby/verbum/backend/pkg/ctl/dictimport/dictparser"
)

func LoadArticles(directory string) (dictparser.Dictionary, error) {
	ids, err := readNewLineSeparatedFile(directory + "/ids.html")
	if err != nil {
		return dictparser.Dictionary{}, fmt.Errorf("read ids file: %w", err)
	}

	headwords, err := readNewLineSeparatedFile(directory + "/headwords.html")
	if err != nil {
		return dictparser.Dictionary{}, fmt.Errorf("read headwords file: %w", err)
	}

	headwordsAlt, err := readNewLineSeparatedFile(directory + "/headwords_alt.html")
	if err != nil {
		return dictparser.Dictionary{}, fmt.Errorf("read headwords alt file: %w", err)
	}

	phrases, err := readNewLineSeparatedFile(directory + "/phrases.html")
	if err != nil {
		return dictparser.Dictionary{}, fmt.Errorf("read phrases file: %w", err)
	}

	bodies, err := readNewLineSeparatedFile(directory + "/bodies.html")
	if err != nil {
		return dictparser.Dictionary{}, fmt.Errorf("read bodies file: %w", err)
	}

	lens := []int{len(ids), len(headwords), len(headwordsAlt), len(phrases), len(bodies)}
	for _, l := range lens[1:] {
		if lens[0] != l {
			return dictparser.Dictionary{}, fmt.Errorf("number of objects in files don't match: %v", lens)
		}
	}

	n := len(ids)
	articles := []dictparser.Article{}
	for i := 0; i < n; i++ {
		a := dictparser.Article{
			ID:           strings.Split(ids[i], "\n")[0],
			Headwords:    strings.Split(headwords[i], "\n"),
			HeadwordsAlt: []string{},
			Phrases:      []string{},
			Body:         bodies[i],
		}
		a.Title = strings.Join(a.Headwords, ", ")

		if headwordsAlt[i] != "-" {
			a.HeadwordsAlt = strings.Split(headwordsAlt[i], "\n")
		}

		if phrases[i] != "-" {
			a.Phrases = strings.Split(phrases[i], "\n")
		}

		articles = append(articles, a)
	}

	return dictparser.Dictionary{Articles: articles, IDsProvided: true}, nil
}

func readNewLineSeparatedFile(filename string) ([]string, error) {
	contentbytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read %s file: %w", filename, err)
	}

	content := string(contentbytes)
	stylesEnd := "</style>\n\n"
	stylesEndPos := strings.LastIndex(content, stylesEnd)
	if stylesEndPos > -1 {
		content = content[stylesEndPos+len(stylesEnd):]
	}

	if content[len(content)-1] == '\n' {
		content = content[:len(content)-1]
	}

	return strings.Split(content, "\n\n"), nil
}
