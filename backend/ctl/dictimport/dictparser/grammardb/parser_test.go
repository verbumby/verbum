package grammardb

import (
	"embed"
	"fmt"
	"os"
	"testing"
)

//go:embed stubdict
var stubdict embed.FS

//go:embed stubids.txt
var stubids []byte

func TestEverything(t *testing.T) {
	tmp, err := os.MkdirTemp("", "grammardb-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	if err := os.WriteFile(tmp+"/ids.txt", stubids, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := os.Mkdir(tmp+"/grammardb", 0o755); err != nil {
		t.Fatal(err)
	}

	stubdictEntries, err := stubdict.ReadDir("stubdict")
	if err != nil {
		t.Fatal(err)
	}

	for _, sde := range stubdictEntries {
		content, err := stubdict.ReadFile("stubdict/" + sde.Name())
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(tmp+"/grammardb/"+sde.Name(), content, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	articlesCh, errCh := ParseDirectory(tmp)

	outHTML, err := os.Create("grammardb.html")
	if err != nil {
		t.Fatal(err)
	}
	defer outHTML.Close()
	_, _ = outHTML.WriteString(`
		<style>
		table, tr, td, th {
			border: 1px solid lightgrey;
			vertical-align: top;
			padding: 0.2rem;
		}
		table {border-collapse: collapse;}
		v-abbr {color:darkgreen;font-weight: normal;}
		hr {border: none; border-top: 1px solid lightgray;}
		</style>
	`)

	for a := range articlesCh {
		fmt.Println(a.ID)
		fmt.Println(a.Headwords)
		fmt.Println(a.HeadwordsAlt)
		fmt.Println(a.Title)
		fmt.Println("--------------------")
		fmt.Println(a.Body)
		fmt.Println("^^^^^^^^^^^^^^^^^^^^")
		if _, err := outHTML.WriteString(a.Body); err != nil {
			t.Fatal(err)
		}
		if _, err := outHTML.WriteString("\n<hr/>\n"); err != nil {
			t.Fatal(err)
		}
	}

	if err := <-errCh; err != nil {
		t.Fatal(err)
	}
}
