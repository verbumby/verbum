package textutil

import (
	"fmt"
	"testing"
)

func TestRomanizeRussian(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{s: "привет", want: "privet"},
		{s: "Арагонская равнина", want: "Aragonskaya ravnina"},
		{s: "абвгдеёжзиклмнопрстуфхцчшщъыьэюя", want: "abvgdeyozhziklmnoprstufkhtschshshchyeyuya"},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case%d", i), func(t *testing.T) {
			if got := RomanizeRussian(tt.s); got != tt.want {
				t.Errorf("RomanizeRussian() = %v, want %v", got, tt.want)
			}
		})
	}
}
