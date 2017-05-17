package token

import (
	"encoding/json"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// FromJSON extract an oauth from a json string
func FromJSON(str string) (*oauth2.Token, error) {
	var token oauth2.Token
	if err := json.Unmarshal([]byte(str), &token); err != nil {
		return nil, errors.Wrap(err, "failed unmarshall json token")
	}
	return &token, nil
}
