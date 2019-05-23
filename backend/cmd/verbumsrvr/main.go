package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/verbumby/verbum/backend/pkg/app"
	"github.com/verbumby/verbum/backend/pkg/chttp"
	"github.com/verbumby/verbum/backend/pkg/handlers"
	"github.com/verbumby/verbum/backend/pkg/tm"
)

func main() {
	if err := app.Bootstrap(); err != nil {
		log.Fatal(err)
	}

	if err := bootstrapServer(); err != nil {
		log.Fatal(err)
	}
}

func bootstrapServer() error {
	templates := []struct {
		name    string
		files   []string
		funcMap template.FuncMap
	}{
		{
			name: "index",
			files: []string{
				"./templates/layout.html",
				"./templates/index-page.html",
				"./templates/search-control.html",
			},
			funcMap: template.FuncMap{},
		},
		{
			name: "search-results",
			files: []string{
				"./templates/layout.html",
				"./templates/search-results-page.html",
				"./templates/search-control.html",
				"./templates/article.html",
			},
			funcMap: template.FuncMap{},
		},
		{
			name: "dictionary",
			files: []string{
				"./templates/layout.html",
				"./templates/dictionary-page.html",
				"./templates/article.html",
				"./templates/pagination.html",
			},
			funcMap: template.FuncMap{},
		},
		{
			name: "article",
			files: []string{
				"./templates/layout.html",
				"./templates/article-page.html",
				"./templates/article.html",
			},
			funcMap: template.FuncMap{},
		},
	}

	for _, t := range templates {
		if err := tm.Compile(t.name, t.files, t.funcMap); err != nil {
			return errors.Wrapf(err, "compile %s template", t.name)
		}
	}

	r := mux.NewRouter()
	statics := http.FileServer(http.Dir("statics"))
	r.HandleFunc("/robots.txt", chttp.MakeHandler(handlers.RobotsTXT))
	r.HandleFunc("/sitemap-index.xml", chttp.MakeHandler(handlers.SitemapIndex))
	r.PathPrefix("/statics/").Handler(http.StripPrefix("/statics/", statics))
	r.HandleFunc("/_suggest", chttp.MakeHandler(handlers.Suggest))
	r.HandleFunc("/{dictionary:[a-z]+}", chttp.MakeHandler(handlers.Dictionary))
	r.HandleFunc("/{dictionary:[a-z]+}/sitemap-{n:[0-9]+}.xml", chttp.MakeHandler(handlers.SitemapOfDictionary))
	r.HandleFunc("/{dictionary:[a-z]+}/{article:[a-zA-Z0-9-]+}", chttp.MakeHandler(handlers.Article))
	r.HandleFunc("/", chttp.MakeHandler(handlers.Index))

	chttp.InitCookieManager()

	if viper.IsSet("http.addr") {
		go func() {
			statics := http.FileServer(http.Dir(viper.GetString("http.acmeChallengeRoot")))
			r := http.NewServeMux()
			r.Handle("/.well-known/acme-challenge/", http.StripPrefix("/.well-known/acme-challenge/", statics))
			r.HandleFunc("/", handlers.ToHTTPS)
			http.ListenAndServe(viper.GetString("http.addr"), r)
		}()
	}

	log.Printf("listening on %s", viper.GetString("https.addr"))
	err := http.ListenAndServeTLS(
		viper.GetString("https.addr"),
		viper.GetString("https.certFile"),
		viper.GetString("https.keyFile"),
		r,
	)
	if err != nil {
		return errors.Wrap(err, "listen and serve")
	}

	return nil
}
