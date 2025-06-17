package dictionary

import "testing"

func Test_wrapAccentedChars(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "case1",
			s:    ">а́<",
			want: `><span class="accent">а́</span><`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wrapAccentedChars(tt.s)
			if got != tt.want {
				t.Errorf("wrapAccentedChars() = %v, want %v", got, tt.want)
			}
		})
	}
}
