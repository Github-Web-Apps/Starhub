package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env"
	"github.com/caarlos0/watchub/datastores/database"
	"github.com/caarlos0/watchub/dto"
	"github.com/caarlos0/watchub/scheduler"
	"github.com/caarlos0/watchub/static"
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
	OauthState   string `env:"OAUTH_STATE"`
	DatabaseURL  string `env:"DATABASE_URL" envDefault:"postgres://localhost:5432/watchub?sslmode=disable"`
}

func main() {
	log.Println("Starting up...")
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
		Scopes:       []string{"user:email", "public_repo"},
		Endpoint:     githuboauth.Endpoint,
	}

	// schedulers
	scheduler := scheduler.New(store)
	scheduler.Start()
	defer scheduler.Stop()

	e := echo.New()
	e.SetRenderer(static.New("static/*.html"))
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", dto.User{})
	})
	e.GET("/login", func(c echo.Context) error {
		url := oauthConf.AuthCodeURL(cfg.OauthState, oauth2.AccessTypeOnline)
		return c.Redirect(http.StatusTemporaryRedirect, url)
	})
	e.GET("/github_callback", func(c echo.Context) error {
		state := c.FormValue("state")
		if state != cfg.OauthState {
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
		if err := store.SaveToken(*u.ID, token); err != nil {
			return err
		}
		if err := store.Schedule(*u.ID, time.Now()); err != nil {
			return err
		}
		return c.Render(http.StatusOK, "index", dto.User{User: *u.Login})
	})
	e.GET("/executions", func(c echo.Context) error {
		executions, err := store.Executions()
		if err != nil {
			return err
		}
		return c.JSON(200, executions)
	})
	e.Run(standard.New(fmt.Sprintf(":%d", cfg.Port)))
}
