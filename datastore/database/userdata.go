package database

import (
	"encoding/json"

	"github.com/caarlos0/watchub/shared/model"
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
func (db *Userdatastore) GetStars(userID int64) (result []model.Star, err error) {
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
func (db *Userdatastore) SaveStars(userID int64, stars []model.Star) error {
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

// FollowerCount returns the amount of followers stored for a given userID
func (db *Userdatastore) FollowerCount(userID int64) (count int, err error) {
	err = db.QueryRow(
		"SELECT array_length(followers, 1) FROM tokens WHERE user_id = $1",
		userID,
	).Scan(&count)
	return
}

const starCountQuery = `
	SELECT SUM( json_array_length( (repo->>'stargazers')::json ) )
	FROM tokens t
	CROSS JOIN json_array_elements(t.stars) repo
	WHERE t.user_id = $1
`

// StarCount returns the amount of stargazers of all the user's repositories
func (db *Userdatastore) StarCount(userID int64) (count int, err error) {
	err = db.QueryRow(starCountQuery, userID).Scan(&count)
	return
}

// RepositoryCount returns the amount of followers stored for a given userID
func (db *Userdatastore) RepositoryCount(userID int64) (count int, err error) {
	err = db.QueryRow(
		"SELECT json_array_length(stars) FROM tokens WHERE user_id = $1",
		userID,
	).Scan(&count)
	return
}
