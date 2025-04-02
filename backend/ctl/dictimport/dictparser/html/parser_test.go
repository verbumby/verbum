package html

import (
	_ "embed"
	"reflect"
	"strings"
	"testing"

	"github.com/verbumby/verbum/backend/dictionary"
)

//go:embed test.html
var testHTML string

var thirdArticleHTML = `<p><strong class="hw">A</strong>` + "\u200b" + `<sup>2</sup> <v-trx>[eɪ]</v-trx> <em>n.</em> «выда́тна» (<em>самая высокая акадэмічная адзнака ў Англіі</em>);</p>
<p class="ms-5"><v-ex><em>He got an A in chemistry.</em> Ён атрымаў «выдатна» па хіміі.</v-ex></p>` + "\n"

func TestParse(t *testing.T) {
	articlesCh, errCh := ParseReader(strings.NewReader(testHTML), dictionary.IndexSettings{})

	expecteds := []struct {
		id     string
		title  string
		hws    []string
		hwsalt []string
	}{
		{
			id:    "A",
			title: "A, a",
			hws:   []string{"A", "a"},
		},
		{
			id:    "A-1",
			title: "A 1",
			hws:   []string{"A"},
		},
		{
			id:    "A-2",
			title: "A 2",
			hws:   []string{"A"},
		},
		{
			id:    "a",
			title: "a",
			hws:   []string{"a"},
		},
		{
			id:    "AA",
			title: "AA",
			hws:   []string{"AA"},
		},
		{
			id:    "aback",
			title: "aback",
			hws:   []string{"aback"},
		},
		{
			id:    "adapter",
			title: "adapter, adaptor",
			hws:   []string{"adapter", "adaptor"},
		},
		{
			id:     "run-2",
			title:  "run 2",
			hws:    []string{"run"},
			hwsalt: []string{"run across", "run around", "run away", "run down", "run into", "run off", "run on", "run out", "run over", "run through", "run up"},
		},
		{
			id:    "parenthetic",
			title: "parenthetic(al)",
			hws:   []string{"parenthetic", "parenthetical"},
		},
		{
			id:    "angina",
			title: "angina (pectoris)",
			hws:   []string{"angina", "angina pectoris"},
		},
	}

	i := 0
	for a := range articlesCh {
		expected := expecteds[i]
		if i == 2 {
			if a.Body != thirdArticleHTML {
				t.Fatal("the body of 3rd article doesn't match")
			}
		}

		if expected.id != a.ID {
			t.Errorf("article %d: ID doesn't match: expected %s, got %s", i, expected.id, a.ID)
		}
		if expected.title != a.Title {
			t.Errorf("article %d: Title doesn't match: expected %s, got %s", i, expected.title, a.Title)
		}
		if !reflect.DeepEqual(expected.hws, a.Headwords) {
			t.Errorf("article %d: Headwords don't match: expected %v, got %v", i, expected.hws, a.Headwords)
		}
		if !reflect.DeepEqual(expected.hwsalt, a.HeadwordsAlt) {
			t.Errorf("article %d: HeadwordsAlt don't match: expected %v, got %v", i, expected.hwsalt, a.HeadwordsAlt)
		}
		i++
	}

	if i != 10 {
		t.Fatalf("expected 10 articles, got %d", i)
	}

	err := <-errCh
	if err != nil {
		t.Fatal(err)
	}
}

func TestUnbalancedParenthesis(t *testing.T) {
	sourceHTML := `<p><strong class="hw">эпіс́омы (ад</strong></p>`
	articlesCh, errCh := ParseReader(strings.NewReader(sourceHTML), dictionary.IndexSettings{})
	for a := range articlesCh {
		t.Fatalf("Expected no articles to be returned, got %v", a)
	}
	err := <-errCh
	if err == nil {
		t.Fatal("Expected unbalanced parenthesis to be detected")
	}
}
