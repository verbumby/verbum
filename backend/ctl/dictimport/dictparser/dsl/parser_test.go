package dsl

import (
	_ "embed"
	"reflect"
	"strings"
	"testing"
)

func TestPrepareHeadwordsForIndexing(t *testing.T) {
	cases := []struct {
		name string
		hws  []string
		want []string
	}{
		{
			name: "curly braces removed with it's inner content",
			hws: []string{
				"abridg{(}e{)}ment",
				"abridg{(e)}ment",
			},
			want: []string{
				"abridgement",
				"abridgment",
			},
		},
		{
			name: "trim whitespaces",
			hws:  []string{" asdf	"},
			want: []string{"asdf"},
		},
		{
			name: "escaped parentheses",
			hws:  []string{"ВКП\\(б\\)"},
			want: []string{"ВКП(б)"},
		},
		{
			name: "{[']}A{[/']}bbau",
			hws:  []string{"{[']}A{[/']}bbau"},
			want: []string{"Abbau"},
		},
	}

	for _, c := range cases {
		actual := prepareHeadwordsForIndexing(c.hws)
		if !reflect.DeepEqual(actual, c.want) {
			t.Errorf("expected %v got %v", c.want, actual)
		}
	}
}

func TestAssembleTitleFromHeadwords(t *testing.T) {
	cases := []struct {
		name string
		hws  []string
		want string
	}{
		{
			name: "duplicate headwords are accounted for",
			hws: []string{
				"abridg{(}e{)}ment",
				"abridg{(e)}ment",
			},
			want: "abridg(e)ment",
		},
		{
			name: "curly braces are removed without it's inner content",
			hws: []string{
				"abridg{(}e{)}ment",
			},
			want: "abridg(e)ment",
		},
		{
			name: "trim whitespaces",
			hws:  []string{" asdf	"},
			want: "asdf",
		},
		{
			name: "escaped parentheses",
			hws:  []string{"ВКП\\(б\\)"},
			want: "ВКП(б)",
		},
	}

	for _, c := range cases {
		actual := assembleTitleFromHeadwords(c.hws)
		if !reflect.DeepEqual(actual, c.want) {
			t.Errorf("expected %v got %v", c.want, actual)
		}
	}
}

//go:embed test.dsl
var testDSL string

func TestParse(t *testing.T) {
	articlesCh, errCh := ParseReader(strings.NewReader(testDSL))

	expecteds := []struct {
		title  string
		hws    []string
		hwsalt []string
	}{
		{
			title: "[']a[/']alglatt",
			hws:   []string{"aalglatt"},
		},
		{
			title: "abridg(e)ment",
			hws:   []string{"abridgement", "abridgment"},
		},
		{
			title: "ВКП(б)",
			hws:   []string{"ВКП(б)"},
		},
		{
			title: "(the) Milky Way",
			hws:   []string{"Milky Way", "the Milky Way"},
		},
		{
			title:  "вярхоўны, Вярхоўны Савет БССР, Вярхоўны Суд БССР",
			hws:    []string{"вярхоўны"},
			hwsalt: []string{"Вярхоўны Савет БССР", "Вярхоўны Суд БССР"},
		},
		{
			title:  "абняць, (як) вокам (зрокам, позіркам) абняць",
			hws:    []string{"абняць"},
			hwsalt: []string{"(як) вокам (зрокам, позіркам) абняць"},
		},
	}

	i := 0
	for a := range articlesCh {
		expected := expecteds[i]

		if i == 0 {
			if a.Body != "[m1][p]a[/p] сл[']і[/']зкі як вуг[']о[/']р; выкр[']у[/']тлівы[/m]\n\n" {
				t.Fatal("the body of 0 article doesn't match")
			}
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

	if i != 6 {
		t.Fatalf("expected 6 articles, got %d", i)
	}

	err := <-errCh
	if err != nil {
		t.Fatal(err)
	}
}
