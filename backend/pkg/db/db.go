package db

import (
	"database/sql"
	"log"

	"github.com/pkg/errors"
	reform "gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/mysql"
)

var (
	// DB reform database handler
	DB *reform.DB
)

// Initialize initializes db connection
func Initialize(connectionString string) error {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return errors.Wrap(err, "open db")
	}
	DB = reform.NewDB(db, mysql.Dialect, reform.NewPrintfLogger(log.Printf))
	return nil
}
