package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	gorillahandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/verbumby/verbum/backend/pkg/app"
	"github.com/verbumby/verbum/backend/pkg/chttp"
	"github.com/verbumby/verbum/backend/pkg/handlers"
	"github.com/verbumby/verbum/backend/pkg/storage"
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
	r := mux.NewRouter()
	r.HandleFunc("/api/dictionaries/{dictionary:[a-z-]+}/articles/{article:[a-zA-Z0-9-]+}", chttp.MakeHandler(handlers.APIArticle, chttp.ContentTypeJSONMiddleware))
	r.HandleFunc("/api/dictionaries/{dictionary:[a-z-]+}/letterfilter", chttp.MakeHandler(handlers.APILetterFilter, chttp.ContentTypeJSONMiddleware))
	r.HandleFunc("/api/dictionaries/{dictionary:[a-z-]+}/articles", chttp.MakeHandler(handlers.APIDictionaryArticles, chttp.ContentTypeJSONMiddleware))
	r.HandleFunc("/api/dictionaries", chttp.MakeHandler(handlers.APIDictionariesList, chttp.ContentTypeJSONMiddleware))
	r.HandleFunc("/api/search", chttp.MakeHandler(handlers.APISearch, chttp.ContentTypeJSONMiddleware))
	r.HandleFunc("/api/suggest", chttp.MakeHandler(handlers.APISuggest, chttp.ContentTypeJSONMiddleware))
	r.PathPrefix("/api/").HandlerFunc(chttp.MakeHandler(handlers.APINotFound))
	imagesServer := http.FileServer(http.Dir(viper.GetString("images.path")))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images", imagesServer))
	r.HandleFunc("/robots.txt", chttp.MakeHandler(handlers.RobotsTXT))
	r.HandleFunc("/sitemap-index.xml", chttp.MakeHandler(handlers.SitemapIndex))
	r.HandleFunc("/{dictionary:[a-z-]+}/sitemap-{n:[0-9]+}.xml", chttp.MakeHandler(handlers.SitemapOfDictionary))
	rpurl := url.URL{Scheme: "http", Host: "localhost:8079"}
	rp := httputil.NewSingleHostReverseProxy(&rpurl)
	r.PathPrefix("/").Handler(rp)

	chttp.InitCookieManager()
	go storage.PruneOldBackups()

	if viper.IsSet("http.addr") {
		go func() {
			statics := http.FileServer(http.Dir(viper.GetString("http.acmeChallengeRoot")))
			r := http.NewServeMux()
			r.Handle("/.well-known/acme-challenge/", http.StripPrefix("/.well-known/acme-challenge/", statics))
			r.HandleFunc("/", handlers.ToHTTPS)
			log.Printf("listening on %s", viper.GetString("http.addr"))
			http.ListenAndServe(viper.GetString("http.addr"), r)
		}()
	}

	rootHandler := gorillahandlers.RecoveryHandler()(
		gorillahandlers.CompressHandler(r),
	)

	privateServer := &http.Server{
		Addr:         viper.GetString("http.private.addr"),
		Handler:      rootHandler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", viper.GetString("http.private.addr"))
		err := privateServer.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Printf("private server listen and serve: %v", err)
		}
	}()

	publicServer := &http.Server{
		Addr:         viper.GetString("https.addr"),
		Handler:      rootHandler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", viper.GetString("https.addr"))
		err := publicServer.ListenAndServeTLS(
			viper.GetString("https.certFile"),
			viper.GetString("https.keyFile"),
		)
		if err != http.ErrServerClosed {
			log.Printf("public server listen and serve tls: %v", err)
		}
	}()

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	log.Println("shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := publicServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("public server shutdown: %v", err)
	}

	if err := privateServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("private server shutdown: %v", err)
	}

	log.Println("see ya!")

	return nil
}
