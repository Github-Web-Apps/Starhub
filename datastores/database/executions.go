package database

import (
	"github.com/caarlos0/watchub/datastores"
	"github.com/jmoiron/sqlx"
)

// Execstore in database
type Execstore struct {
	*sqlx.DB
}

// NewExecstore datastore
func NewExecstore(db *sqlx.DB) *Execstore {
	return &Execstore{db}
}

// Executions get the executions that should be made
func (db *Execstore) Executions() ([]datastores.Execution, error) {
	var executions []datastores.Execution
	return executions, db.Select(
		&executions, "SELECT user_id, token FROM tokens WHERE next <= now()",
	)
}
