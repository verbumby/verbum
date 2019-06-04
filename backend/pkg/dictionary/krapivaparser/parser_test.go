package krapivaparser

import "testing"

func TestKrapivaParser(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
		wantErr bool
	}{
		{
			name:    "escape sequence 1",
			content: "\\[\\][c darkblue][m1][']а[/']",
			want:    `[]<span style="color: darkblue"><p class="ml-0"><span style="color: brown;">а́</span>`,
			wantErr: false,
		},
		{
			name:    "case1",
			content: "[m1][b][c darkblue]непрыгляднасць,[/c][/b][c brown] ‑і, [p]ж.[/p][/c]\n[m2]\\[Вася\\] Уласцівасць непрыгляднага; непрывабнасць.\n",
			want:    `<p class="ml-0"><strong><span style="color: darkblue">непрыгляднасць,</span></strong><span style="color: brown"> ‑і, ж.</span></p><p class="ml-2">[Вася] Уласцівасць непрыгляднага; непрывабнасць.</p>`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := KrapivaParser(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("KrapivaParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("KrapivaParser() = %v, want %v", got, tt.want)
			}
		})
	}
}
