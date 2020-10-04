package app

import (
	"fmt"

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
		return fmt.Errorf("read in config: %w", err)
	}

	if err := article.Migrate(); err != nil {
		return fmt.Errorf("article migrate: %w", err)
	}

	if err := elasticIndexTemplatesMigrate(); err != nil {
		return fmt.Errorf("elastic index templates migrate: %w", err)
	}

	return nil
}

func elasticIndexTemplatesMigrate() error {
	respbody := &struct {
		AccessLog struct {
			Version int
		} `json:"access-log"`
	}{}
	err := storage.Get("/_template/access-log?filter_path=*.version", respbody)

	if err != nil || respbody.AccessLog.Version != 1 {
		err = storage.Put("/_template/access-log", map[string]interface{}{
			"version":        1,
			"index_patterns": []string{"access-log-*"},
			"settings": map[string]interface{}{
				"number_of_shards":   1,
				"number_of_replicas": 0,
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
		if err != nil {
			return fmt.Errorf("create access log index template: %w", err)
		}
	}

	return nil
}
