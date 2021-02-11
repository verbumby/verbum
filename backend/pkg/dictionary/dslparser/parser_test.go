package dslparser

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
			want:    `[]<span style="color: darkblue"><p class="ml-0"><v-accent>а́</v-accent>`,
			wantErr: false,
		},
		{
			name:    "case1",
			content: "[m1][b][c darkblue]непрыгляднасць,[/c][/b][c brown] ‑і, [p]ж.[/p][/c]\n[m2]\\[Вася\\] Уласцівасць непрыгляднага; непрывабнасць.\n",
			want:    `<p class="ml-0"><strong><span style="color: darkblue">непрыгляднасць,</span></strong><span style="color: brown"> ‑і, ж.</span></p><p class="ml-2">[Вася] Уласцівасць непрыгляднага; непрывабнасць.</p>`,
			wantErr: false,
		},
		// {
		// 	name: "case2",
		// 	content: "	[m1][b]Том:[/b] 2, [b]старонка:[/b] 26.[/m]\n[s]02-026_0079_\[no_name\].jpg[/s]\n",

		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DSLParser(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("DSLParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DSLParser() = %v, want %v", got, tt.want)
			}
		})
	}
}
