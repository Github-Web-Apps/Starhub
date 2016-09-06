package database

import (
	"encoding/json"

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

func (db *Tokenstore) SaveToken(userID int, token *oauth2.Token) error {
	strToken, err := tokenToJSON(token)
	if err != nil {
		return err
	}
	previousToken, err := db.GetUserToken(userID)
	if previousToken != nil && err == nil {
		_, err := db.Exec(
			"UPDATE tokens SET token = $2, updated_at = now() WHERE user_id = $1",
			userID,
			strToken,
		)
		return err
	}
	_, err = db.Exec(
		"INSERT INTO tokens(user_id, token) VALUES($1, $2)", userID, strToken,
	)
	return err
}

func (db *Tokenstore) GetUserToken(userID int) (*oauth2.Token, error) {
	var token string
	err := db.Get(
		&token, "SELECT token FROM tokens WHERE user_id = $1", userID,
	)
	if err != nil {
		return nil, err
	}
	return tokenFromJSON(token)
}

func tokenToJSON(token *oauth2.Token) (string, error) {
	d, err := json.Marshal(token)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func tokenFromJSON(jsonStr string) (*oauth2.Token, error) {
	var token oauth2.Token
	if err := json.Unmarshal([]byte(jsonStr), &token); err != nil {
		return nil, err
	}
	return &token, nil
}
