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

func (db *Userdatastore) GetFollowers(userID int64) ([]int64, error) {
	var ids []sql.NullInt64
	var result []int64
	err := db.QueryRowx(
		"SELECT followers FROM tokens WHERE user_id = $1",
		userID,
	).Scan(pq.Array(&ids))
	if err != nil {
		return result, err
	}
	for _, id := range ids {
		if id.Valid {
			result = append(result, id.Int64)
		}
	}
	return result, nil
}

func (db *Userdatastore) SaveFollowers(userID int64, followers []int64) error {
	_, err := db.Exec(
		"UPDATE tokens SET followers = $2 WHERE user_id = $1",
		userID,
		pq.Array(followers),
	)
	return err
}
