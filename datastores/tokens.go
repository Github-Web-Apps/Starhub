package datastores

import (
	"time"

	"golang.org/x/oauth2"
)

type Tokenstore interface {
	SaveToken(userID int, token *oauth2.Token) error
	Schedule(userID int, date time.Time) error
}
