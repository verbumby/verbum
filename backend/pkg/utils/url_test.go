package utils

import "testing"

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
		{name: "case2", args: args{s: "абвгдеёжзіийклмнопрстуўфхцчшщъыьэюя"}, want: "аbvgdyeyozhziiyklmnoprstuufkhtschshSHCH''y'eyuya"},
		{name: "case3", args: args{s: "АБВГДЕЁЖЗІИЙКЛМНОПРСТУЎФХЦЧШЩЪЫЬЭЮЯ"}, want: "АBVGDYEYOZHZIIYKLMNOPRSTUUFKHTSCHSHSHCH''Y'EYUYA"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Latinize(tt.args.s); got != tt.want {
				t.Errorf("Latinize() = %v, want %v", got, tt.want)
			}
		})
	}
}
