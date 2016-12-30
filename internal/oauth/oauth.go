package oauth

import (
	"errors"
	"net/http"

	"context"

	"github.com/caarlos0/watchub/internal/config"
	"github.com/caarlos0/watchub/internal/datastores"
	"github.com/caarlos0/watchub/internal/dto"
	"github.com/google/go-github/github"
	"github.com/labstack/echo"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

const applicationsURL = "https://github.com/settings/connections/applications/"

// Oauth info
type Oauth struct {
	config *oauth2.Config
	store  datastores.Datastore
	state  string
}

// New oauth
func New(store datastores.Datastore, config config.Config) *Oauth {
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

// Client for a given token
func (o *Oauth) Client(token *oauth2.Token) *github.Client {
	return github.NewClient(o.config.Client(context.Background(), token))
}

// Mount the routes as a group of a given echo instance
func (o *Oauth) Mount(e *echo.Echo) *echo.Group {
	login := e.Group("login")

	login.GET("", func(c echo.Context) error {
		url := o.config.AuthCodeURL(o.state, oauth2.AccessTypeOnline)
		return c.Redirect(http.StatusTemporaryRedirect, url)
	})

	login.GET("/callback", func(c echo.Context) error {
		state := c.FormValue("state")
		if state != o.state {
			return errors.New("invalid state")
		}
		code := c.FormValue("code")
		token, err := o.config.Exchange(context.Background(), code)
		if err != nil {
			return err
		}
		client := github.NewClient(o.config.Client(context.Background(), token))
		u, _, err := client.Users.Get("")
		if err != nil {
			return err
		}
		if err := o.store.SaveToken(int64(*u.ID), token); err != nil {
			return err
		}
		return c.Render(http.StatusOK, "index", dto.User{
			User: *u.Login,
			ChangeSubscriptionURL: applicationsURL + o.config.ClientID,
		})
	})

	return login
}
