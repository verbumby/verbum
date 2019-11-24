package textutil

import (
	"fmt"
	"testing"
)

func TestRomanizeBelarusian(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{s: "Аршанскі", want: "Arshanski"},
		{s: "Ельск", want: "Jelsk"},
		{s: "Бабаедава", want: "Babajedava"},
		{s: "Лепель", want: "Liepiel"},
		{s: "Тлумачальны слоўнік беларускай мовы", want: "Tlumachalny slounik bielaruskaj movy"},
		{s: "коп'епадобны", want: "kopjepadobny"},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case%d", i), func(t *testing.T) {
			if got := RomanizeBelarusian(tt.s); got != tt.want {
				t.Errorf("RomanizeBelarusian() = %v, want %v", got, tt.want)
			}
		})
	}
}
