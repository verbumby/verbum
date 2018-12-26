package article

import (
	"reflect"
	"testing"
)

func Test_RvblrParse(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name       string
		args       args
		wantHws    []string
		wantHwsalt []string
		wantErr    bool
	}{
		{
			name: "case1",
			args: args{
				content: `<v-hw>аазіс</v-hw>

				*назоўнік | мужчынскі род*

				Месца ў пустыні, дзе ёсць вада і расліннасць.

				* Зялёны а. у моры пяскоў.
				* Культурны а. (пераноснае значэнне).

				|| прыметнік: <v-hw>аазісны</v-hw>. Аазісная расліннасць.`,
			},
			wantHws:    []string{"аазіс"},
			wantHwsalt: []string{"аазісны"},
			wantErr:    false,
		}, {
			name: "case2",
			args: args{
				content: `<v-hw>аб'ект</v-hw>

				*назоўнік | мужчынскі род*

				1. У філасофіі: тое, што існуе па-за намі і незалежна ад нашай свддомасці, навакольны свет, матэрыяльная рэчаіснасць (спецыяльны тэрмін).

				2. З'ява, прадмет, на які накіравана чыя-н. дзейнасць.
					* А. навуковага даследавання.
					* А. назірання.

				3. Прадпрыемства, будоўля, установа, а таксама ўсё тое, што з'яўляецца месцам якой-н. дзейнасці.
					* А. будаўніцтва.
					* Пускавы а.

				4. У граматыцы: семантычныя катэгорыі са значэннем таго, на каго (што) накіравана дзеянне.

				|| прыметнік: <v-hw>аб'ектны</v-hw> і <v-hw>аб'ектавы</v-hw>.`,
			},
			wantHws:    []string{"аб'ект"},
			wantHwsalt: []string{"аб'ектны", "аб'ектавы"},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHws, gotHwsalt, err := RvblrParse(tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("rvblrParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHws, tt.wantHws) {
				t.Errorf("rvblrParse() gotHws = %v, want %v", gotHws, tt.wantHws)
			}
			if !reflect.DeepEqual(gotHwsalt, tt.wantHwsalt) {
				t.Errorf("rvblrParse() gotHwsalt = %v, want %v", gotHwsalt, tt.wantHwsalt)
			}
		})
	}
}
