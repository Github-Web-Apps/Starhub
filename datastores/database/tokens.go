package database

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/oauth2"
)

// Tokenstore in database
type Tokenstore struct {
	*sqlx.DB
}

// NewTokenstore datastore
func NewTokenstore(db *sqlx.DB) *Tokenstore {
	return &Tokenstore{db}
}

const insertTokenStm = `
	INSERT INTO tokens(user_id, token, next)
	VALUES($1, $2, now())
	ON CONFLICT(user_id)
		DO UPDATE SET token = $2, updated_at = now(), next = now();
`

func (db *Tokenstore) SaveToken(userID int64, token *oauth2.Token) error {
	strToken, err := tokenToJSON(token)
	if err != nil {
		return err
	}
	_, err = db.Exec(insertTokenStm, userID, strToken)
	return err
}

func (db *Tokenstore) Schedule(userID int64, date time.Time) error {
	_, err := db.Exec(
		"UPDATE tokens SET next = $2, updated_at = now() WHERE user_id = $1",
		userID,
		date,
	)
	return err
}

func tokenToJSON(token *oauth2.Token) (string, error) {
	d, err := json.Marshal(token)
	if err != nil {
		return "", err
	}
	return string(d), nil
}
