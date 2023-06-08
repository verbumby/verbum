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

	gorillahandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/verbumby/verbum/backend/chttp"
	"github.com/verbumby/verbum/backend/ctl"
	"github.com/verbumby/verbum/backend/ctl/dictimport"
	"github.com/verbumby/verbum/backend/dictionary"
	"github.com/verbumby/verbum/backend/handlers"
	"github.com/verbumby/verbum/backend/serverrender"
	"github.com/verbumby/verbum/backend/storage"
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
		ctl.ReindexCommand(),
		ctl.MigrateSlugs(),
		ctl.MigrateRvblrWrongHeadwords(),
		ctl.MigrateStardictFixPrefixCase(),
		ctl.MigrateModifiedAt(), ctl.MigrateTitleToContentCommand(),
		ctl.MigrateZeroWidthSpaceBeforeSupCommand(),
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

	viper.SetDefault("cookie.name", "vadm")
	viper.SetDefault("cookie.nameState", "vadm-state")
	viper.SetDefault("cookie.maxAge", 604800)

	viper.SetDefault("oauth.endpointToken", "https://www.googleapis.com/oauth2/v4/token")
	viper.SetDefault("oauth.endpointUserinfo", "https://www.googleapis.com/oauth2/v3/userinfo")
	viper.SetDefault("oauth.endpointAuth", "https://accounts.google.com/o/oauth2/v2/auth")

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
	r := mux.NewRouter()
	r.HandleFunc("/api/dictionaries/{dictionary:[a-z-]+}/articles/{article:[a-zA-Z0-9-]+}", chttp.MakeHandler(handlers.APIArticle, chttp.ContentTypeJSONMiddleware))
	r.HandleFunc("/api/dictionaries/{dictionary:[a-z-]+}/letterfilter", chttp.MakeHandler(handlers.APILetterFilter, chttp.ContentTypeJSONMiddleware))
	r.HandleFunc("/api/dictionaries/{dictionary:[a-z-]+}/articles", chttp.MakeHandler(handlers.APIDictionaryArticles, chttp.ContentTypeJSONMiddleware))
	r.HandleFunc("/api/dictionaries/{dictionary:[a-z-]+}/abbrevs", chttp.MakeHandler(handlers.APIDictionaryAbbrevs, chttp.ContentTypeJSONMiddleware))
	r.HandleFunc("/api/dictionaries", chttp.MakeHandler(handlers.APIDictionariesList, chttp.ContentTypeJSONMiddleware))
	r.HandleFunc("/api/search", chttp.MakeHandler(handlers.APISearch, chttp.ContentTypeJSONMiddleware))
	r.HandleFunc("/api/suggest", chttp.MakeHandler(handlers.APISuggest, chttp.ContentTypeJSONMiddleware))
	r.PathPrefix("/api/").HandlerFunc(chttp.MakeHandler(handlers.APINotFound))
	imagesServer := http.FileServer(http.Dir(viper.GetString("images.path")))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images", imagesServer))
	staticsServer := http.FileServer(http.FS(frontend.DistPublic))
	staticsHander := http.StripPrefix("/statics", staticsServer)
	r.PathPrefix("/statics/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		staticsHander.ServeHTTP(w, r)
	})
	r.HandleFunc("/robots.txt", chttp.MakeHandler(handlers.RobotsTXT))
	r.HandleFunc("/sitemap-index.xml", chttp.MakeHandler(handlers.SitemapIndex))
	r.HandleFunc("/{dictionary:[a-z-]+}/sitemap-{n:[0-9]+}.xml", chttp.MakeHandler(handlers.SitemapOfDictionary))
	serverRenderer, err := serverrender.New(r)
	if err != nil {
		return fmt.Errorf("create server renderer: %w", err)
	}
	r.PathPrefix("/").Handler(chttp.MakeHandler(serverRenderer.ServeHTTP))

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
