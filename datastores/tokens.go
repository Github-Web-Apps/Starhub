package datastores

import "golang.org/x/oauth2"

type Tokenstore interface {
	Save(userID int, token *oauth2.Token) error
	ForUser(userID int) (*oauth2.Token, error)
}
