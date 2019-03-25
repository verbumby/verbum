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

func TestSlugify(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{s: "some___url *** !@#!@#//slug-", want: "some-url-slug"},
		{s: "Tlumachalny slounik bielaruskaj movy", want: "tlumachalny-slounik-bielaruskaj-movy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Slugify(tt.s); got != tt.want {
				t.Errorf("Slugify() = %v, want %v", got, tt.want)
			}
		})
	}
}
