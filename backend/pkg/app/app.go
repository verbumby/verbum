package app

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/verbumby/verbum/backend/pkg/db"
	"github.com/verbumby/verbum/backend/pkg/fts"
)

// Bootstrap bootstraps the application
func Bootstrap() error {
	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("fts.host", "localhost")

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

	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "read in config")
	}

	// TODO: parametrize db connection
	dbConnString := fmt.Sprintf("root@tcp(%s:3306)/verbum", viper.GetString("db.host"))
	if err := db.Initialize(dbConnString); err != nil {
		return errors.Wrap(err, "db initalize")
	}

	sphinxConnString := fmt.Sprintf("tcp(%s:9306)/?interpolateParams=true", viper.GetString("fts.host"))
	if err := fts.Initialize(sphinxConnString); err != nil {
		return errors.Wrap(err, "fts initialize")
	}

	return nil
}
