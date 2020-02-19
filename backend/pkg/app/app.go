package app

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/verbumby/verbum/backend/pkg/article"
	"github.com/verbumby/verbum/backend/pkg/storage"
)

// Bootstrap bootstraps the application
func Bootstrap() error {
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

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/usr/local/share/verbum")

	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "read in config")
	}

	if err := article.Migrate(); err != nil {
		return errors.Wrap(err, "article migrate")
	}

	if err := elasticIndexTemplatesMigrate(); err != nil {
		return errors.Wrap(err, "elastic index templates migrate")
	}

	return nil
}

func elasticIndexTemplatesMigrate() error {
	err := storage.Query(http.MethodHead, "/_template/access-log", nil, nil)
	if err == nil {
		return nil
	}

	err = storage.Put("/_template/access-log", map[string]interface{}{
		"index_patterns": []string{"access-log-*"},
		"settings": map[string]interface{}{
			"number_of_shards": 1,
		},
		"mappings": map[string]interface{}{
			"_doc": map[string]interface{}{
				"properties": map[string]interface{}{
					"TS":        map[string]interface{}{"type": "date"},
					"Path":      map[string]interface{}{"type": "keyword"},
					"Query":     map[string]interface{}{"type": "keyword"},
					"IP":        map[string]interface{}{"type": "keyword"},
					"UserAgent": map[string]interface{}{"type": "keyword"},
					"Referer":   map[string]interface{}{"type": "keyword"},
				},
			},
		},
	}, nil)
	return errors.Wrap(err, "create access log index template")
}
