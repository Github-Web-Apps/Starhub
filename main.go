package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/datastores/database"
	"github.com/caarlos0/watchub/dto"
	"github.com/caarlos0/watchub/oauth"
	"github.com/caarlos0/watchub/scheduler"
	"github.com/caarlos0/watchub/static"
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

	// schedulers
	scheduler := scheduler.New(store)
	scheduler.Start()
	defer scheduler.Stop()

	// routes
	e := echo.New()
	e.SetRenderer(static.New("static/*.html"))
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", dto.User{})
	})
	oauth.Mount(e, store, config)

	e.Run(standard.New(fmt.Sprintf(":%d", config.Port)))
}
