package fts

import (
	"database/sql"
	"log"

	"github.com/pkg/errors"
	reform "gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/mysql"
)

var (
	// Sphinx is a reform connection adapter to sphinxsearch
	Sphinx *reform.DB
)

// Initialize initializes sphinx connection
func Initialize(connectionString string) error {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return errors.Wrap(err, "open sphinx")
	}
	Sphinx = reform.NewDB(db, mysql.Dialect, reform.NewPrintfLogger(log.Printf))
	return nil
}
