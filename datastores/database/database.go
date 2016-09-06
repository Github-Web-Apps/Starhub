package database

import (
	"database/sql"
	"log"

	"github.com/caarlos0/watchub/datastores"
	"github.com/jmoiron/sqlx"
)

// Connect creates a connection pool to the database
func Connect(url string) *sql.DB {
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

// NewDatastore returns a new Datastore
func NewDatastore(db *sql.DB) datastores.Datastore {
	dbx := sqlx.NewDb(db, "postgres")
	return struct {
		*Tokenstore
	}{
		NewTokenstore(dbx),
	}
}
