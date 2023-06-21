package dslparser

import "testing"

func TestKrapivaParser(t *testing.T) {
	tests := []struct {
		name    string
		content string
		opts    []Option
		want    string
		wantErr bool
	}{
		{
			name:    "escape sequence 1",
			content: "\\[\\][c darkblue][m1][']а[/']",
			want:    `[]<span style="color: darkblue"><p class="ms-0"><v-accent>а́</v-accent>`,
			wantErr: false,
		},
		{
			name:    "case1",
			content: "[m1][b][c darkblue]непрыгляднасць,[/c][/b][c brown] ‑і, [p]ж.[/p][/c]\n[m2]\\[Вася\\] Уласцівасць непрыгляднага; непрывабнасць.\n",
			want:    `<p class="ms-0"><strong><span style="color: darkblue">непрыгляднасць,</span></strong><span style="color: brown"> ‑і, <v-abbr>ж.</v-abbr></span></p><p class="ms-2">[Вася] Уласцівасць непрыгляднага; непрывабнасць.</p>`,
			wantErr: false,
		},
		{
			name:    "case2",
			content: "[m1][b]Том:[/b] 2, [b]старонка:[/b] 26.[/m]\n[s]02-026_0079_\\[no_name\\].jpg[/s]\n",
			opts:    []Option{GlobalStore("dictID", "dict-id-1")},
			want:    `<p class="ms-0"><strong>Том:</strong> 2, <strong>старонка:</strong> 26.</p><img src="/images/dict-id-1/02/02-026_0079_%5Bno_name%5D.jpg" alt="02-026_0079_%5Bno_name%5D.jpg"/></p>`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse("article", []byte(tt.content), tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
