package database

import (
	"database/sql"

	"github.com/apex/log"
	"github.com/caarlos0/watchub/internal/datastores"
	"github.com/jmoiron/sqlx"
)

// Connect creates a connection pool to the database
func Connect(url string) *sql.DB {
	var log = log.WithField("url", url)
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.WithError(err).Fatal("Failed to open connection to database")
	}
	if err := db.Ping(); err != nil {
		log.WithError(err).Fatal("Failed to ping database")
	}
	return db
}

// NewDatastore returns a new Datastore
func NewDatastore(db *sql.DB) datastores.Datastore {
	dbx := sqlx.NewDb(db, "postgres")
	return struct {
		*Tokenstore
		*Execstore
		*Userdatastore
	}{
		NewTokenstore(dbx),
		NewExecstore(dbx),
		NewUserdatastore(dbx),
	}
}
