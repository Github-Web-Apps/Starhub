package main

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/caarlos0/watchub/internal/config"
	"github.com/caarlos0/watchub/internal/datastores/database"
	"github.com/caarlos0/watchub/internal/dto"
	"github.com/caarlos0/watchub/internal/oauth"
	"github.com/caarlos0/watchub/internal/pages"
	"github.com/caarlos0/watchub/internal/scheduler"
	"github.com/gorilla/mux"
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
	defer func() { _ = db.Close() }()
	store := database.NewDatastore(db)

	// oauth
	oauth := oauth.New(store, config)

	// schedulers
	scheduler := scheduler.New(config, store, oauth)
	scheduler.Start()
	defer scheduler.Stop()

	// routes
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.Methods("GET").Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pages.Render(w, "index", dto.IndexData{})
	})
	r.Methods("GET").Path("/donate").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pages.Render(w, "donate", dto.IndexData{})
	})
	r.Methods("GET").Path("/support").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pages.Render(w, "support", dto.IndexData{})
	})

	// mount oauth routes
	oauth.Mount(r)

	// RUN!
	var server = &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:" + config.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
