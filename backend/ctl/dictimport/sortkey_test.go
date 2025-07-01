package dictimport

import "testing"

func TestSortKey(t *testing.T) {
	pairs := [][2]string{
		{"а", "аб"},
		{"а", "б"},
		{"аба", "аб’едкі"},
	}

	for _, p := range pairs {
		a, b := p[0], p[1]
		ak, bk := createSortKey(a), createSortKey(b)
		if ak >= bk {
			t.Errorf("Invalid order: %s (%s) must be smaller than %s (%s)", a, ak, b, bk)
		}
	}
}
