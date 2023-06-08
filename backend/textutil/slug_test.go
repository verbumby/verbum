package textutil

import "testing"

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
