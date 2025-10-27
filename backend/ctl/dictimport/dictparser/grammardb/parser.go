package grammardb

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"path"
	"regexp"
	"slices"
	"strings"

	xmlparser "github.com/tamerh/xml-stream-parser"
	"github.com/verbumby/verbum/backend/ctl/dictimport/dictparser"
)

//go:embed template.html
var bodyTemplateSource string

func lemma(hw string) string {
	hw = strings.ReplaceAll(hw, "+", "\u0301")
	hw = strings.ReplaceAll(hw, `'`, `’`)
	return hw
}

func hasForm(data *renderDataT, tag string) bool {
	tagRe := regexp.MustCompile(`^` + tag + `$`)
	for _, form := range data.Forms {
		if tagRe.MatchString(form.Tag) {
			return true
		}
	}
	return false
}

func countForms(data *renderDataT, forms ...string) int {
	result := 0
	for _, f := range forms {
		if hasForm(data, f) {
			result++
		}
	}
	return result
}

func form(data *renderDataT, tag string) (template.HTML, error) {
	tagRe := regexp.MustCompile(`^` + tag + `$`)
	forms := []*renderDataFormT{}
	for _, form := range data.Forms {
		if tagRe.MatchString(form.Tag) {
			form.Visited = true
			forms = append(forms, form)
		}
	}

	if len(forms) == 0 {
		return "-", nil
		// return "", fmt.Errorf("couldn't find any forms with tag %s", tag)
	}

	words := []string{}
	for _, form := range forms {
		word := form.Lemma
		if form.Option != "" {
			opt := ""
			switch form.Option {
			case "inanim":
				opt = "неадуш."
			case "anim":
				opt = "адуш."
			default:
				return "", fmt.Errorf("unknown form option: %s", form.Option)
			}
			opt = fmt.Sprintf(" (<v-abbr>%s</v-abbr>)", opt)
			word += opt
		}
		words = append(words, word)
	}

	return template.HTML(strings.Join(words, "<br/>")), nil
}

func grammarCases() []map[string]string {
	return []map[string]string{
		{"Tag": "N", "Abbr": "Н."},
		{"Tag": "G", "Abbr": "Р."},
		{"Tag": "D", "Abbr": "Д."},
		{"Tag": "A", "Abbr": "В."},
		{"Tag": "I", "Abbr": "Т."},
		{"Tag": "L", "Abbr": "М."},
	}
}

func grammarCasesNoun() []map[string]string {
	result := grammarCases()
	result = append(result, map[string]string{"Tag": "V", "Abbr": "Кл."})
	return result
}

func list(a ...any) []any {
	return a
}

type renderDataFormT struct {
	Tag     string
	Lemma   string
	Option  string
	Visited bool
}

type renderDataT struct {
	Tag           string
	Lemma         string
	Meaning       string
	OtherVariants []string
	Forms         []*renderDataFormT
	Sources       []string
}

func ParseDirectory(dirname string) (chan dictparser.Article, chan error) {
	articlesCh := make(chan dictparser.Article, 64)
	errCh := make(chan error)

	go func() {
		retErr := func(format string, a ...any) {
			close(articlesCh)
			errCh <- fmt.Errorf(format, a...)
			close(errCh)
		}

		funcMap := template.FuncMap{
			"list":              list,
			"hasForm":           hasForm,
			"countForms":        countForms,
			"form":              form,
			"hasPrefix":         strings.HasPrefix,
			"hasSuffix":         strings.HasSuffix,
			"grammarCategories": grammarCategories,
			"grammarCases":      grammarCases,
			"grammarCasesNoun":  grammarCasesNoun,
		}
		bodyTemplate, err := template.New("body").Funcs(funcMap).Parse(bodyTemplateSource)
		if err != nil {
			retErr("parse body template: %w", err)
			return
		}

		idsf, err := os.Open(dirname + "/ids.txt")
		if err != nil {
			retErr("open ids.txt: %w", err)
			return
		}
		defer idsf.Close()
		idssc := bufio.NewScanner(idsf)

		files, err := os.ReadDir(dirname)
		if err != nil {
			retErr("read dir %s: %w", dirname, err)
			return
		}

		n := 0
		body := &bytes.Buffer{}
		for _, file := range files {
			if path.Ext(file.Name()) != ".xml" {
				continue
			}

			filename := dirname + "/" + file.Name()
			fmt.Println(filename)

			f, err := os.Open(filename)
			if err != nil {
				retErr("open %s: %w", filename, err)
				return
			}
			defer f.Close()

			parser := xmlparser.NewXMLParser(bufio.NewReaderSize(f, 65_536), "Paradigm")
			for paradigmXML := range parser.Stream() {
				if paradigmXML.Err != nil {
					retErr("parse paradigm: %w", paradigmXML.Err)
					return
				}

				for _, variantXML := range paradigmXML.Childs["Variant"] {
					n++
					body.Reset()

					vid := paradigmXML.Attrs["pdgId"] + "-" + variantXML.Attrs["id"]

					tag := paradigmXML.Attrs["tag"]
					if v, ok := variantXML.Attrs["tag"]; ok {
						tag += v
					}

					if tag[0] == 'K' {
						continue
					}

					if !idssc.Scan() {
						retErr("ids stream has ended preliminary")
						return
					}
					idref := strings.Split(idssc.Text(), "\t")

					l := strings.ReplaceAll(variantXML.Attrs["lemma"], "+", "")
					if idref[0] != vid || idref[1] != tag || idref[2] != l {
						retErr("idref and the source don't match: %v vs. %s %s %s", idref, vid, tag, l)
						return
					}

					data := &renderDataT{
						Tag:     tag,
						Lemma:   lemma(variantXML.Attrs["lemma"]),
						Meaning: paradigmXML.Attrs["meaning"],
						OtherVariants: func() (result []string) {
							for _, otherVariant := range paradigmXML.Childs["Variant"] {
								if otherVariant.Attrs["lemma"] != variantXML.Attrs["lemma"] {
									result = append(result, lemma(otherVariant.Attrs["lemma"]))
								}
							}
							return
						}(),
						Forms: func() (result []*renderDataFormT) {
							for _, formXML := range variantXML.Childs["Form"] {
								result = append(result, &renderDataFormT{
									Tag:    formXML.Attrs["tag"],
									Lemma:  lemma(formXML.InnerText),
									Option: formXML.Attrs["options"],
								})
							}
							return
						}(),
						Sources: func() []string {
							sourceLists := []string{variantXML.Attrs["slouniki"]}
							for _, formXML := range variantXML.Childs["Form"] {
								sourceLists = append(sourceLists, formXML.Attrs["slouniki"])
							}

							set := map[string]bool{}
							for _, list := range sourceLists {
								for _, s := range strings.Split(list, ",") {
									s = strings.TrimSpace(s)
									if s == "" {
										continue
									}
									if strings.Contains(s, ":") {
										sp := strings.Split(s, ":")
										s = sp[0]
									}
									set[s] = true
								}
							}

							result := []string{}
							for k := range set {
								result = append(result, k)
							}
							slices.Sort(result)
							return result
						}(),
					}
					if err := bodyTemplate.Execute(body, data); err != nil {
						retErr("render body article body: %w", err)
						return
					}
					hws := headwords(variantXML)
					outid := hws[0]
					if idref[3] != "" {
						outid += "-" + idref[3]
					}
					articlesCh <- dictparser.Article{
						ID:           outid,
						Title:        hws[0],
						Headwords:    hws,
						HeadwordsAlt: []string{},
						Body:         body.String(),
					}

					for _, form := range data.Forms {
						if !form.Visited && form.Tag != "" && form.Tag != "0" && form.Tag != "1" && form.Tag != "XXX" && !strings.HasSuffix(form.Tag, "HX") {
							retErr("form %v was not visited", form)
							return
						}
					}
				}
			}
		}
		fmt.Println(n)
		close(articlesCh)
		close(errCh)
	}()

	return articlesCh, errCh
}

func headwords(variantXML xmlparser.XMLElement) []string {
	lemma := variantXML.Attrs["lemma"]
	lemma = strings.ReplaceAll(lemma, "+", "")
	lemma = strings.ReplaceAll(lemma, "'", "’")
	return []string{lemma}
}

func grammarCategories(tag string) string {
	result := []string{}

	switch tag[0] {
	case 'N':
		result = append(result, "назоўнік")

		switch tag[1] {
		case 'C':
			result = append(result, "агульны")
		case 'P':
			result = append(result, "уласны")
		}

		switch tag[2] {
		case 'A':
			result = append(result, "адушаўлёны")
		case 'I':
			result = append(result, "неадушаўлёны")
		}

		switch tag[3] {
		case 'P':
			result = append(result, "асабовы")
		case 'I':
			result = append(result, "неасабовы")
		}

		// tag[4]

		switch tag[5] {
		case 'M':
			result = append(result, "мужчынскі род")
		case 'F':
			result = append(result, "жаночы род")
		case 'N':
			result = append(result, "ніякі род")
		case 'C':
			result = append(result, "агульны")
		case 'S':
			result = append(result, "субстантываваны")
		case 'U':
			result = append(result, "субстантываваны множналікавы")
		case 'P':
			result = append(result, "множны лік")
		}

		switch tag[6] {
		case '1':
			result = append(result, "1 скланенне")
		case '2':
			result = append(result, "2 скланенне")
		case '3':
			result = append(result, "3 скланенне")
		case '0':
			result = append(result, "нескланяльны")
		case '4':
			result = append(result, "рознаскланяльны")
		case '6':
			result = append(result, "змешанае скланенне")
		case '5':
			result = append(result, "ад’ектыўнае скланенне")
		case '7':
			result = append(result, "множналікавы")
		}

	case 'A':
		result = append(result, "прыметнік")

		switch tag[1] {
		case 'Q':
			result = append(result, "якасны")
		case 'R':
			result = append(result, "адносны")
		case 'P':
			result = append(result, "прыналежны")
		case '0':
			result = append(result, "нескланяльны")
		}

		if tag[1] != '0' {
			switch tag[2] {
			case 'P':
			// result = append(result, "станоўчая cтупень параўнання")
			case 'C':
				result = append(result, "вышэйшая cтупень параўнання")
			case 'S':
				result = append(result, "найвышэйшая cтупень параўнання")
			}
		}
	case 'C':
		result = append(result, "злучнік")
		switch tag[1] {
		case 'S':
			result = append(result, "падпарадкавальны")
		case 'K':
			result = append(result, "злучальны")
		}
	case 'E':
		result = append(result, "часціца")
	case 'I':
		result = append(result, "прыназоўнік")
	case 'M':
		result = append(result, "лічэбнік")

		if tag[1] == '0' {
			result = append(result, "нескланяльны")
		}

		switch tag[2] {
		case 'C':
			result = append(result, "колькасны")
		case 'O':
			result = append(result, "парадкавы")
		case 'K':
			result = append(result, "зборны")
		case 'F':
			result = append(result, "дробавы")
		}

	case 'P':
		result = append(result, "дзеепрыметнік")

		switch tag[1] {
		case 'A':
			result = append(result, "незалежны стан")
		case 'P':
			result = append(result, "залежны стан")
		}

		switch tag[2] {
		case 'R':
			result = append(result, "цяперашні час")
		case 'P':
			result = append(result, "прошлы час")
		}

		switch tag[3] {
		case 'P':
			result = append(result, "закончанае трыванне")
		case 'M':
			result = append(result, "незакончанае трыванне")
		}
	case 'R':
		result = append(result, "прыслоўе")

		switch tag[1] {
		case 'N':
			result = append(result, "утворана ад назоўніка")
		case 'A':
			result = append(result, "утворана ад прыметніка")
		case 'M':
			result = append(result, "утворана ад лічэбніка")
		case 'S':
			result = append(result, "утворана ад займенніка")
		case 'G':
			result = append(result, "утворана ад дзеепрыслоўя")
		case 'V':
			result = append(result, "утворана ад дзеяслова")
		case 'E':
			result = append(result, "утворана ад часціцы")
		case 'I':
			result = append(result, "утворана ад прыназоўніка")
		}

	case 'S':
		result = append(result, "займеннік")
		switch tag[2] {
		case 'P':
			result = append(result, "асабовы")
		case 'R':
			result = append(result, "зваротны")
		case 'S':
			result = append(result, "прыналежны")
		case 'D':
			result = append(result, "указальны")
		case 'E':
			result = append(result, "азначальны")
		case 'L':
			result = append(result, "пытальна–адносны")
		case 'N':
			result = append(result, "адмоўны")
		case 'F':
			result = append(result, "няпэўны")
		}

		switch tag[3] {
		case '1':
			result = append(result, "1-я асоба")
		case '2':
			result = append(result, "2-я асоба")
		case '3':
			result = append(result, "3-я асоба")
		case '0':
			result = append(result, "безасабовы")
		}

	case 'V':
		result = append(result, "дзеяслоў")
		switch tag[1] {
		case 'T':
			result = append(result, "пераходны")
		case 'I':
			result = append(result, "непераходны")
		case 'D':
			result = append(result, "пераходны/непераходны")
		}

		switch tag[2] {
		case 'P':
			result = append(result, "закончанае трыванне")
		case 'M':
			result = append(result, "незакончанае трыванне")
		}

		switch tag[3] {
		case 'R':
			result = append(result, "зваротны")
		case 'N':
			result = append(result, "незваротны")
		}

		switch tag[4] {
		case '1':
			result = append(result, "1-е спражэнне")
		case '2':
			result = append(result, "2-е спражэнне")
		case '3':
			result = append(result, "рознаспрагальны")
		}

	case 'W':
		result = append(result, "прэдыкатыў")
	case 'Y':
		result = append(result, "выклічнік")
	case 'Z':
		result = append(result, "пабочнае слова")
	}

	return strings.Join(result, ", ")
}
