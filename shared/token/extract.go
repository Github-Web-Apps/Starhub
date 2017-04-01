package token

import (
	"encoding/json"

	"golang.org/x/oauth2"
)

func FromJSON(str string) (*oauth2.Token, error) {
	var token oauth2.Token
	if err := json.Unmarshal([]byte(str), &token); err != nil {
		return nil, err
	}
	return &token, nil
}
