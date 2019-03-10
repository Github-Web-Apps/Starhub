package controllers

import (
	"net/http"
	"time"

	"github.com/Intika-Web-Apps/Starhub-Notifier/config"
	"github.com/Intika-Web-Apps/Starhub-Notifier/datastore"
	"github.com/Intika-Web-Apps/Starhub-Notifier/shared/pages"
	"github.com/gorilla/sessions"
)

// Schedule ctrl
type Schedule struct {
	Base
	store datastore.Datastore
}

// NewSchedule ctrl
func NewSchedule(
	config config.Config,
	session sessions.Store,
	store datastore.Datastore,
) *Schedule {
	return &Schedule{
		Base: Base{
			config:  config,
			session: session,
		},
		store: store,
	}
}

// Handler handles /Schedule
func (ctrl *Schedule) Handler(w http.ResponseWriter, r *http.Request) {
	session, _ := ctrl.session.Get(r, ctrl.config.SessionName)
	id, _ := session.Values["user_id"].(int)
	if session.IsNew || id == 0 {
		http.Error(w, "not logged in", http.StatusForbidden)
		return
	}
	if err := ctrl.store.Schedule(int64(id), time.Now()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pages.Render(w, "scheduled", ctrl.sessionData(w, r))
}
