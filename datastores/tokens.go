package datastores

import (
	"time"

	"golang.org/x/oauth2"
)

type Tokenstore interface {
	SaveToken(userID int64, token *oauth2.Token) error
	Schedule(userID int64, date time.Time) error
}
