package oauth

import (
	"context"

	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/datastore"
	"github.com/caarlos0/watchub/shared/token"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

// Oauth info
type Oauth struct {
	config *oauth2.Config
	store  datastore.Datastore
	state  string
}

// New oauth
func New(store datastore.Datastore, config config.Config) *Oauth {
	return &Oauth{
		config: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			Scopes:       []string{"user:email,public_repo"},
			Endpoint:     githuboauth.Endpoint,
		},
		store: store,
		state: config.OauthState,
	}
}

// ClientFrom for a given string token
func (o *Oauth) ClientFrom(ctx context.Context, tokenStr string) (*github.Client, error) {
	token, err := token.FromJSON(tokenStr)
	if err != nil {
		return nil, err
	}
	return github.NewClient(o.config.Client(ctx, token)), err
}
