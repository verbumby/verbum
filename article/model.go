package article

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

//go:generate reform

// Article represents an article
//
// reform:articles
type Article struct {
	ID      int32  `reform:"id,pk" json:"id"`
	Content string `reform:"content" json:"content"`
}

// BeforeInsert implements reform.BeforeInserter
func (a *Article) BeforeInsert() error {
	return a.BeforeSave()
}

// BeforeUpdate implements reform.BeforeUpdater
func (a *Article) BeforeUpdate() error {
	return a.BeforeSave()
}

// BeforeSave updates sphinx index/indices
func (a *Article) BeforeSave() error {

	queue := []html.Token{}
	z := html.NewTokenizer(strings.NewReader(a.Content))
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			err := z.Err()
			if err != io.EOF {
				return errors.Wrap(err, "tokenize html")
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
				return fmt.Errorf("no starting tag found for closing tag %s", t.Type)
			}
			if err := processTokenQueue(queue[startingTagIDX:]); err != nil {
				return errors.Wrap(err, "processing tags sequence")
			}
		}

	}
	content := ""
	for _, t := range queue {
		content += t.String()
	}
	a.Content = content
	return nil
}

func processTokenQueue(q []html.Token) error {
	if len(q) < 2 {
		return fmt.Errorf("expected at least two (openning and closing) tokens")
	}

	if q[0].Type != html.StartTagToken {
		return fmt.Errorf("expected starting token type to be %s, got %s", html.StartTagToken, q[0].Type)
	}

	if q[len(q)-1].Type != html.EndTagToken {
		return fmt.Errorf("expected ending token type to be %s, got %s", html.EndTagToken, q[len(q)-1].Type)
	}

	// check whether start and end tokens match
	if q[0].Data != q[len(q)-1].Data {
		return fmt.Errorf("closing tag %s doesn't match opening tag %s", q[len(q)-1].Data, q[0].Data)
	}

	return nil
}
