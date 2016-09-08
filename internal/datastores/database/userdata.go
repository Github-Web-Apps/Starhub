package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Userdatastore in database
type Userdatastore struct {
	*sqlx.DB
}

// NewUserdatastore datastore
func NewUserdatastore(db *sqlx.DB) *Userdatastore {
	return &Userdatastore{db}
}

// GetFollowers of a given userID
func (db *Userdatastore) GetFollowers(userID int64) ([]string, error) {
	var logins []sql.NullString
	var result []string
	err := db.QueryRow(
		"SELECT followers FROM tokens WHERE user_id = $1",
		userID,
	).Scan(pq.Array(&logins))
	if err != nil {
		return result, err
	}
	for _, id := range logins {
		if id.Valid {
			result = append(result, id.String)
		}
	}
	return result, nil
}

// SaveFollowers for a given userID
func (db *Userdatastore) SaveFollowers(userID int64, followers []string) error {
	_, err := db.Exec(
		"UPDATE tokens SET followers = $2 WHERE user_id = $1",
		userID,
		pq.Array(followers),
	)
	return err
}
