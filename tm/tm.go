package tm

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

var templates = map[string]*template.Template{}

// Compile compiles a named templates from files
func Compile(name string, files []string) error {
	t, err := template.New(name).Funcs(template.FuncMap{
		"staticURL": staticURL,
	}).ParseFiles(files...)
	if err != nil {
		return errors.Wrap(err, "parse files")
	}

	templates[name] = t
	return nil
}

// Render renders the name template to w writer
func Render(name string, w io.Writer, data interface{}) error {
	t, ok := templates[name]
	if !ok {
		return fmt.Errorf("no such template with %s name", name)
	}
	if err := t.ExecuteTemplate(w, "layout", data); err != nil {
		return errors.Wrap(err, "execute template")
	}

	return nil
}

func staticURL(file string) string {
	info, err := os.Stat("." + file)
	if err != nil {
		log.Fatalln(err)
	}
	return file + "?" + strconv.FormatInt(info.ModTime().Unix(), 10)
}
