package dsl

import (
	"reflect"
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
