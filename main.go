package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	rice "github.com/GeertJohan/go.rice"
	"github.com/caarlos0/env"
	"github.com/caarlos0/watch/datastores/database"
	"github.com/google/go-github/github"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

type config struct {
	Port         int    `env:"PORT" envDefault:"3000"`
	ClientID     string `env:"GITHUB_CLIENT_ID"`
	ClientSecret string `env:"GITHUB_CLIENT_SECRET"`
	DatabaseURL  string `env:"DATABASE_URL" envDefault:"postgres://localhost:5432/watchub?sslmode=disable"`
}

func main() {
	cfg := config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalln(err)
	}
	db := database.Connect(cfg.DatabaseURL)
	defer db.Close()
	store := database.NewDatastore(db)

	oauthConf := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Scopes:       []string{"user:email", "repo"},
		Endpoint:     githuboauth.Endpoint,
	}
	oauthStateString := "thisshouldberandom"

	e := echo.New()
	assetHandler := http.FileServer(rice.MustFindBox("static").HTTPBox())
	e.GET("/", standard.WrapHandler(assetHandler))
	e.GET("/static/*", standard.WrapHandler(http.StripPrefix("/static/", assetHandler)))
	e.GET("/login", func(c echo.Context) error {
		url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
		return c.Redirect(http.StatusTemporaryRedirect, url)
	})
	e.GET("/github_callback", func(c echo.Context) error {
		state := c.FormValue("state")
		if state != oauthStateString {
			return errors.New("Invalid state!")
		}
		code := c.FormValue("code")
		token, err := oauthConf.Exchange(oauth2.NoContext, code)
		if err != nil {
			return err
		}
		fmt.Println("Save token here:", token)
		oauthClient := oauthConf.Client(oauth2.NoContext, token)
		client := github.NewClient(oauthClient)
		u, _, err := client.Users.Get("")
		if err != nil {
			return err
		}
		if err := store.Save(*u.ID, token); err != nil {
			return err
		}
		return c.String(200, "Hello, "+*u.Login+"!")
	})
	e.Run(standard.New(":3000"))
}
