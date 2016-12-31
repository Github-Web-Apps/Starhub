package main

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/caarlos0/watchub/internal/config"
	"github.com/caarlos0/watchub/internal/datastores/database"
	"github.com/caarlos0/watchub/internal/dto"
	"github.com/caarlos0/watchub/internal/oauth"
	"github.com/caarlos0/watchub/internal/scheduler"
	"github.com/caarlos0/watchub/internal/template"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("Starting up...")

	// config
	config, err := config.Get()
	if err != nil {
		log.Panicln(err)
	}

	// datastores
	db := database.Connect(config.DatabaseURL)
	defer db.Close()
	store := database.NewDatastore(db)

	// oauth
	oauth := oauth.New(store, config)

	// schedulers
	scheduler := scheduler.New(config, store, oauth)
	scheduler.Start()
	defer scheduler.Stop()

	// routes
	e := echo.New()
	e.SetRenderer(template.New("static/*.html"))
	e.Static("/static", "static")
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", dto.User{})
	})
	e.GET("/donate", func(c echo.Context) error {
		return c.Render(http.StatusOK, "donate", nil)
	})

	// mount oauth routes
	oauth.Mount(e)

	// RUN!
	e.Run(standard.New(fmt.Sprintf(":%d", config.Port)))
}
