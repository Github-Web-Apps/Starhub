package main

import (
	"net/http"
	"time"

	_ "net/http/pprof"

	_ "github.com/lib/pq"

	"github.com/apex/httplog"
	"github.com/apex/log"
	"github.com/apex/log/handlers/logfmt"

	"github.com/Intika-Web-Apps/Watchub-Mirror/config"
	"github.com/Intika-Web-Apps/Watchub-Mirror/controllers"
	"github.com/Intika-Web-Apps/Watchub-Mirror/datastore/database"
	"github.com/Intika-Web-Apps/Watchub-Mirror/oauth"
	"github.com/Intika-Web-Apps/Watchub-Mirror/scheduler"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.SetHandler(logfmt.Default)
	log.SetLevel(log.InfoLevel)
	log.Info("starting up...")

	var config = config.Get()
	var db = database.Connect(config.DatabaseURL)
	defer func() { _ = db.Close() }()
	var store = database.NewDatastore(db)

	// oauth
	var session = sessions.NewCookieStore([]byte(config.SessionSecret))
	var oauth = oauth.New(config)
	var loginCtrl = controllers.NewLogin(config, session, oauth, store)

	// schedulers
	var sch = scheduler.New(config, store, oauth, session)
	sch.Start()
	defer sch.Stop()

	// routes
	var mux = mux.NewRouter()
	mux.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)
	mux.Methods(http.MethodGet).Path("/").HandlerFunc(
		controllers.NewIndex(config, session, store).Handler,
	)
	mux.Methods(http.MethodGet).Path("/donate").HandlerFunc(
		controllers.NewDonate(config, session).Handler,
	)
	mux.Methods(http.MethodGet).Path("/contact").HandlerFunc(
		controllers.NewContact(config, session).Handler,
	)
	mux.Methods(http.MethodGet).Path("/schedule").HandlerFunc(
		controllers.NewSchedule(config, session, store).Handler,
	)
	mux.Methods(http.MethodGet).Path("/login").HandlerFunc(
		loginCtrl.Handler,
	)
	mux.Methods(http.MethodGet).Path("/login/callback").HandlerFunc(
		loginCtrl.CallbackHandler,
	)
	mux.Path("/logout").HandlerFunc(
		controllers.NewLogout(config, session).Handler,
	)

	// prometheus stuff
	prometheus.MustRegister(scheduler.TimeGauge)
	prometheus.MustRegister(scheduler.ErrorGauge)
	mux.Handle("/metrics", promhttp.Handler())

	mux.PathPrefix("/debug").Handler(http.DefaultServeMux)

	var handler = context.ClearHandler(
		httplog.New(
			mux,
		),
	)

	var server = &http.Server{
		Handler:      handler,
		Addr:         ":" + config.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.WithField("addr", server.Addr).Info("started")
	if err := server.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("failed to start up server")
	}
}
