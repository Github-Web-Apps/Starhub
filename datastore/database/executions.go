package database

import (
	"github.com/caarlos0/watchub/config/model"
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

const executionsStmQuery = `
	UPDATE tokens
	SET next = now() + interval '1 day', updated_at = now()
	WHERE next <= now()
	RETURNING user_id, token
`

// Executions get the executions that should be made
func (db *Execstore) Executions() (executions []model.Execution, err error) {
	return executions, db.Select(&executions, executionsStmQuery)
}
