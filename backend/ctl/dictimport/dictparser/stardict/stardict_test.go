package stardict

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

//go:embed test.dict
var testDict string

func TestParseReader(t *testing.T) {
	articleCh, errCh := LoadArticles(strings.NewReader(testDict))

	expecteds := []struct {
		title string
		hws   []string
		body  string
	}{
		{
			title: "-ка",
			hws:   []string{"-ка"},
			body:  "<k>-ка</k>\nчастица ану, (после гласных) жен. жа, чаще не переводится совсем ну-ка, покажи — ану, пакажы посмотри-ка на него — паглядзі на яго дайте-ка пройти — дайце (ж) прайсці а, давай, а заеду-ка я, в самом деле, к нему — а заеду (давай заеду) я, сапраўды, да яго",
		},
		{
			title: "-либо",
			hws:   []string{"-либо"},
			body:  "<k>-либо</k>\n-небудзь, когда-либо — калі-небудзь где-либо — дзе-небудзь",
		},
		{
			title: "АСУ",
			hws:   []string{"АСУ"},
			body:  "<k>АСУ</k>\nАСК",
		},
	}

	i := 0
	for a := range articleCh {
		expected := expecteds[i]
		if a.Title != expected.title {
			t.Errorf("title doesn't match: expected %s, got %s", expected.title, a.Title)
		}
		if a.Body != expected.body {
			j1, _ := json.Marshal(a.Body)
			fmt.Println(string(j1))
			j2, _ := json.Marshal(expected.body)
			fmt.Println(string(j2))
			t.Errorf("body doesn't match: expected %s, got %s", expected.body, a.Body)
		}
		if !reflect.DeepEqual(a.Headwords, expected.hws) {
			t.Errorf("headwords don't match: expected %v, got %v", expected.hws, a.Headwords)
		}
		i++
	}

	if i != 3 {
		t.Errorf("expected %d articles in total, got %d", 3, i)
	}

	err := <-errCh
	if err != nil {
		t.Fatal(err)
	}
}
