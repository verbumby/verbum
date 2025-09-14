package dictimport

import "testing"

func TestCalcIDBase(t *testing.T) {
	tests := []struct {
		id           string
		expectedBase string
		expectedNo   int
	}{
		{id: "asdf", expectedBase: "asdf", expectedNo: -1},
		{id: "ВАСІЛЕВІЧЫ-2", expectedBase: "ВАСІЛЕВІЧЫ", expectedNo: 2},
	}

	for _, tc := range tests {
		actualBase, actualNo := calcIDBase(tc.id)
		if actualBase != tc.expectedBase {
			t.Errorf("returned base %s, expected %s", actualBase, tc.expectedBase)
		}
		if actualNo != tc.expectedNo {
			t.Errorf("returned no %d, expected %d", actualNo, tc.expectedNo)
		}
	}
}
