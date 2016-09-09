package database

import (
	"encoding/json"

	"github.com/caarlos0/watchub/internal/datastores"
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
	var logins []string
	return logins, db.QueryRow(
		"SELECT followers FROM tokens WHERE user_id = $1",
		userID,
	).Scan(pq.Array(&logins))
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

// GetStars of a given userID
func (db *Userdatastore) GetStars(userID int64) ([]datastores.Star, error) {
	var result []datastores.Star
	var stars json.RawMessage
	if err := db.QueryRow(
		"SELECT stars FROM tokens WHERE user_id = $1",
		userID,
	).Scan(&stars); err != nil {
		return result, err
	}
	return result, json.Unmarshal(stars, &result)
}

// SaveStars for a given userID
func (db *Userdatastore) SaveStars(userID int64, stars []datastores.Star) error {
	data, err := json.Marshal(stars)
	if err != nil {
		return err
	}
	_, err = db.Exec(
		"UPDATE tokens SET stars = $2 WHERE user_id = $1",
		userID,
		data,
	)
	return err
}
