package utils

import (
	"testing"
)

func TestLatinize(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "case1", args: args{s: "слоўнік"}, want: "slounik"},
		{name: "case2", args: args{s: "абвгдеёжзіийклмнопрстуўфхцчшщъыьэюя"}, want: "abvgdyeyozhziiyklmnoprstuufkhtschshshch''y'eyuya"},
		{name: "case3", args: args{s: "АБВГДЕЁЖЗІИЙКЛМНОПРСТУЎФХЦЧШЩЪЫЬЭЮЯ"}, want: "ABVGDYEYOZHZIIYKLMNOPRSTUUFKHTSCHSHSHCH''Y'EYUYA"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Latinize(tt.args.s); got != tt.want {
				t.Errorf("Latinize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlugify(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "case1", args: args{s: "слоўнік"}, want: "slounik"},
		{name: "case2", args: args{s: "абвгдеёжзіийклмнопрстуўфхцчшщъыьэюя"}, want: "abvgdyeyozhziiyklmnoprstuufkhtschshshch''y'eyuya"},
		{name: "case3", args: args{s: "АБВГДЕЁЖЗІИЙКЛМНОПРСТУЎФХЦЧШЩЪЫЬЭЮЯ"}, want: "abvgdyeyozhziiyklmnoprstuufkhtschshshch''y'eyuya"},
		{
			name: "case4",
			args: args{s: "Тлумачальны слоўнік беларускай мовы (rv-blr.com)"},
			want: "tlumachal'ny-slounik-byelaruskay-movy-(rv-blr.com)",
		},
		{
			name: "case5",
			args: args{s: "Тлумачальны Слоўнік па Інфарматыцы, М. І. Савіцкі, 2014"},
			want: "tlumachal'ny-slounik-pa-infarmatytsy,-m.-i.-savitski,-2014",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Slugify(tt.args.s); got != tt.want {
				t.Errorf("Slugify() = %v, want %v", got, tt.want)
			}
		})
	}
}
