package dictionary

import "testing"

func TestRenderAbbrevs(t *testing.T) {
	abbrevs := &Abbrevs{cache: map[string]*Abbrev{}}
	abbrevs.cache["v"] = &Abbrev{
		Keys:  []string{"v"},
		Value: "verb",
	}

	cases := []struct {
		content  string
		expected string
	}{
		{
			content:  `asdf <v-abbr>v</v-abbr> asdf`,
			expected: `asdf <v-abbr data-bs-toggle="tooltip" data-bs-title="verb" tabindex="0">v</v-abbr> asdf`,
		},
		{
			content:  `asdf <v-abbr class="source">v</v-abbr> asdf`,
			expected: `asdf <v-abbr data-bs-toggle="tooltip" data-bs-title="verb" tabindex="0" class="source">v</v-abbr> asdf`,
		},
	}

	for i, c := range cases {
		actual := renderAbbrevs(c.content, abbrevs)
		if actual != c.expected {
			t.Errorf("render abbrevs case %d expected `%s`, got `%s`", i, c.expected, actual)
		}
	}
}
