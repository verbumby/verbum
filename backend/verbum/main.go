package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/verbumby/verbum/backend/chttp"
	"github.com/verbumby/verbum/backend/ctl"
	"github.com/verbumby/verbum/backend/ctl/dictimport"
	"github.com/verbumby/verbum/backend/dictionary"
	"github.com/verbumby/verbum/backend/handlers"
	"github.com/verbumby/verbum/backend/serverrender"
	"github.com/verbumby/verbum/frontend"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatal(err)
	}

	if err := dictionary.InitDictionaries(); err != nil {
		log.Fatal(err)
	}

	rootCmd := &cobra.Command{
		Use:   "verbum",
		Short: "verbum",
		Long:  "verbum",
	}

	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "serve",
			Short: "Start the http(s) servers",
			RunE:  bootstrapServer,
		},
		dictimport.Command(),
		ctl.CleanupCommand(),
		ctl.ExportCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func initConfig() error {
	viper.SetDefault("https.addr", ":8443")
	viper.SetDefault("https.certFile", "cert.pem")
	viper.SetDefault("https.keyFile", "key.pem")
	viper.SetDefault("https.canonicalAddr", "https://localhost:8443")

	viper.SetDefault("elastic.addr", "http://localhost:9200")

	viper.SetDefault("images.path", "./images")

	viper.SetDefault("dicts.repo.path", "../slouniki")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/usr/local/share/verbum")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read in config: %w", err)
	}

	return nil
}

func bootstrapServer(cmd *cobra.Command, args []string) error {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.HandleFunc("/api/dictionaries/{dictionary:[a-z0-9-]+}/articles/{article}", chttp.MakeHandler(handlers.APIArticle, chttp.ContentTypeJSONMiddleware))
	r.HandleFunc("/api/dictionaries/{dictionary:[a-z0-9-]+}/letterfilter", chttp.MakeHandler(handlers.APILetterFilter, chttp.ContentTypeJSONMiddleware))
	r.HandleFunc("/api/dictionaries/{dictionary:[a-z0-9-]+}/articles", chttp.MakeHandler(handlers.APIDictionaryArticles, chttp.ContentTypeJSONMiddleware))
	r.HandleFunc("/api/dictionaries/{dictionary:[a-z0-9-]+}/abbrevs", chttp.MakeHandler(handlers.APIDictionaryAbbrevs, chttp.ContentTypeJSONMiddleware))
	r.HandleFunc("/api/dictionaries", chttp.MakeHandler(handlers.APIDictionariesList, chttp.ContentTypeJSONMiddleware))
	r.HandleFunc("/api/search", chttp.MakeHandler(handlers.APISearch, chttp.ContentTypeJSONMiddleware))
	r.HandleFunc("/api/suggest", chttp.MakeHandler(handlers.APISuggest, chttp.ContentTypeJSONMiddleware))
	r.Mount("/api/", chttp.MakeHandler(handlers.APINotFound))
	imagesServer := http.FileServer(http.Dir(viper.GetString("images.path")))
	r.Mount("/images/", http.StripPrefix("/images", imagesServer))
	staticsServer := http.FileServer(http.FS(frontend.DistPublic))
	staticsHander := http.StripPrefix("/statics", staticsServer)
	r.Mount("/statics/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		staticsHander.ServeHTTP(w, r)
	}))
	r.HandleFunc("/robots.txt", chttp.MakeHandler(handlers.RobotsTXT))
	r.HandleFunc("/sitemap-index.xml", chttp.MakeHandler(handlers.SitemapIndex))
	r.HandleFunc("/{dictionary:[a-z0-9-]+}/sitemap-{n:[0-9]+}.xml", chttp.MakeHandler(handlers.SitemapOfDictionary))
	serverRenderer, err := serverrender.New(r)
	if err != nil {
		return fmt.Errorf("create server renderer: %w", err)
	}
	r.Mount("/", chttp.MakeHandler(serverRenderer.ServeHTTP))

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

	publicServer := &http.Server{
		Addr:         viper.GetString("https.addr"),
		Handler:      r,
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

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	log.Println("shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := publicServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("public server shutdown: %v", err)
	}

	log.Println("see ya!")

	return nil
}
