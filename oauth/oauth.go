package oauth

import (
	"errors"
	"net/http"

	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/datastores"
	"github.com/caarlos0/watchub/dto"
	"github.com/google/go-github/github"
	"github.com/labstack/echo"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

func Mount(
	e *echo.Echo, store datastores.Datastore, config config.Config,
) *echo.Group {
	login := e.Group("login")

	oauthConf := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       []string{"user:email", "public_repo"},
		Endpoint:     githuboauth.Endpoint,
	}

	login.GET("", func(c echo.Context) error {
		url := oauthConf.AuthCodeURL(config.OauthState, oauth2.AccessTypeOnline)
		return c.Redirect(http.StatusTemporaryRedirect, url)
	})

	login.GET("/callback", func(c echo.Context) error {
		state := c.FormValue("state")
		if state != config.OauthState {
			return errors.New("Invalid state!")
		}
		code := c.FormValue("code")
		token, err := oauthConf.Exchange(oauth2.NoContext, code)
		if err != nil {
			return err
		}
		client := github.NewClient(oauthConf.Client(oauth2.NoContext, token))
		u, _, err := client.Users.Get("")
		if err != nil {
			return err
		}
		if err := store.SaveToken(*u.ID, token); err != nil {
			return err
		}
		return c.Render(http.StatusOK, "index", dto.User{User: *u.Login})
	})

	return login
}
