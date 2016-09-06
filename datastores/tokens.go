package datastores

import "golang.org/x/oauth2"

type Tokenstore interface {
	SaveToken(userID int, token *oauth2.Token) error
	GetUserToken(userID int) (*oauth2.Token, error)
}
