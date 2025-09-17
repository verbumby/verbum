package dictimport

import (
	"fmt"
	"html/template"
	"testing"

	"github.com/verbumby/verbum/backend/dictionary"
)

func TestIDGen(t *testing.T) {
	type testcase struct {
		dict     dictionary.Dictionary
		args     [][2]string
		expected []string // TODO: expand to test returned errors
	}

	tcs := []testcase{}
	tcs = append(tcs, testcase{
		dict:     stubdict{slugifier: "none", indexSettings: dictionary.IndexSettings{DictProvidesIDs: true}},
		args:     [][2]string{{"", "ВАСІЛЕВІЧЫ"}, {"", "ВАСІЛЕВІЧЫ"}, {"", "ВАСІЛЕВІЧЫ-2"}},
		expected: []string{"ВАСІЛЕВІЧЫ", "ВАСІЛЕВІЧЫ-2", "ВАСІЛЕВІЧЫ-2-1"},
	})

	for tcn, tc := range tcs {
		t.Run(fmt.Sprintf("tc_%d", tcn), func(t *testing.T) {
			ig := NewIDGen(tc.dict)
			if len(tc.args) != len(tc.expected) {
				t.Fatalf("length of args does not match length of expected returns")
			}

			for i := 0; i < len(tc.args); i++ {
				args := tc.args[i]
				expected := tc.expected[i]
				actual, _ := ig.Gen(args[0], args[1])
				if actual != expected {
					t.Fatalf("expected %s, but got %s", expected, actual)
				}
			}
		})
	}
}

type stubdict struct {
	dictionary.Common
	slugifier     string
	indexSettings dictionary.IndexSettings
}

func (d stubdict) Slugifier() string {
	return d.slugifier
}

func (d stubdict) IndexSettings() dictionary.IndexSettings {
	return d.indexSettings
}

func (d stubdict) ToHTML(content string) template.HTML {
	return ""
}
