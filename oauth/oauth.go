package oauth

import (
	"context"

	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/shared/token"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

// Oauth info
type Oauth struct {
	config *oauth2.Config
	state  string
}

// New oauth
func New(
	config config.Config,
) *Oauth {
	return &Oauth{
		config: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			Scopes:       []string{"user:email"},
            //Scopes:       []string{"user:email,public_repo"}, //old scopes
			Endpoint:     githuboauth.Endpoint,
		},
		state: config.OauthState,
	}
}

// Client for a given token
func (o *Oauth) Client(ctx context.Context, token *oauth2.Token) *github.Client {
	return github.NewClient(o.config.Client(ctx, token))
}

// ClientFrom for a given string token
func (o *Oauth) ClientFrom(ctx context.Context, tokenStr string) (*github.Client, error) {
	token, err := token.FromJSON(tokenStr)
	if err != nil {
		return nil, err
	}
	return o.Client(ctx, token), nil
}

// AuthCodeURL URL to OAuth 2.0 provider's consent page
func (o *Oauth) AuthCodeURL() string {
	return o.config.AuthCodeURL(o.state, oauth2.AccessTypeOnline)
}

// IsStateValid true if state is valid
func (o *Oauth) IsStateValid(state string) bool {
	return o.state == state
}

// Exchange oauth code
func (o *Oauth) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return o.config.Exchange(ctx, code)
}
