package main

import (
	"net/http"
	"time"

	"github.com/apex/httplog"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/datastore/database"
	"github.com/caarlos0/watchub/oauth"
	"github.com/caarlos0/watchub/scheduler"
	"github.com/caarlos0/watchub/shared/dto"
	"github.com/caarlos0/watchub/shared/pages"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	log.SetHandler(text.Default)
	log.SetLevel(log.InfoLevel)
	log.Info("starting up...")

	var config = config.Get()
	var db = database.Connect(config.DatabaseURL)
	defer func() { _ = db.Close() }()
	var store = database.NewDatastore(db)

	// oauth
	var oauth = oauth.New(store, config)

	// schedulers
	var scheduler = scheduler.New(config, store, oauth)
	scheduler.Start()
	defer scheduler.Stop()

	// routes
	var r = mux.NewRouter()
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

	var loginMux = r.Methods("GET").PathPrefix("/login").Subrouter()
	loginMux.Path("").HandlerFunc(oauth.LoginHandler())
	loginMux.Path("/callback").HandlerFunc(oauth.LoginCallbackHandler())

	// RUN!
	var server = &http.Server{
		Handler:      httplog.New(handlers.CompressHandler(r)),
		Addr:         ":" + config.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.WithField("port", config.Port).Info("started")
	log.WithError(server.ListenAndServe()).Error("failed to start up server")
}
