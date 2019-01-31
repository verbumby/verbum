package app

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
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

	return nil
}
