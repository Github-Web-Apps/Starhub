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
	"github.com/caarlos0/watchub/shared/pages"
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
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
	var session = sessions.NewCookieStore([]byte(config.SessionSecret))
	var oauth = oauth.New(store, session, config)

	// schedulers
	var scheduler = scheduler.New(config, store, oauth, session)
	scheduler.Start()
	defer scheduler.Stop()

	// routes
	var mux = mux.NewRouter()
	mux.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Methods("GET").Path("/").
		HandlerFunc(pages.New(config, session).IndexHandler)
	mux.Methods("GET").Path("/donate").
		HandlerFunc(pages.New(config, session).DonateHandler)
	mux.Methods("GET").Path("/support").
		HandlerFunc(pages.New(config, session).SupportHandler)

	var loginMux = mux.Methods("GET").PathPrefix("/login").Subrouter()
	loginMux.Path("").HandlerFunc(oauth.LoginHandler())
	loginMux.Path("/callback").HandlerFunc(oauth.LoginCallbackHandler())
	mux.Path("/check").HandlerFunc(scheduler.ScheduleHandler())

	var handler = context.ClearHandler(
		httplog.New(
			handlers.CompressHandler(
				mux,
			),
		),
	)
	var server = &http.Server{
		Handler:      handler,
		Addr:         ":" + config.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.WithField("addr", server.Addr).Info("started")
	log.WithError(server.ListenAndServe()).Error("failed to start up server")
}
