package textutil

import "testing"

func TestSortKey(t *testing.T) {
	pairs := [][2]string{
		{"а", "аб"},
		{"а", "б"},
		{"аба", "аб’едкі"},
		{"г", "ґ"},
		{"ґ", "д"},
	}

	for _, p := range pairs {
		a, b := p[0], p[1]
		ak, bk := CreateSortKey(a), CreateSortKey(b)
		if ak >= bk {
			t.Errorf("Invalid order: %s (%s) must be smaller than %s (%s)", a, ak, b, bk)
		}
	}
}
