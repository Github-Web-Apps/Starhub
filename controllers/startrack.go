package controllers

import (
	"net/http"

	"github.com/Intika-Web-Apps/Starhub-Notifier/config"
	"github.com/Intika-Web-Apps/Starhub-Notifier/shared/pages"
	"github.com/gorilla/sessions"
)

// Startrack ctrl
type Startrack struct {
	Base
}

// NewStartrack ctrl
func NewStartrack(
	config config.Config,
	session sessions.Store,
) *Startrack {
	return &Startrack{
		Base: Base{
			config:  config,
			session: session,
		},
	}
}

// Handler handles /Startrack
func (ctrl *Startrack) Handler(w http.ResponseWriter, r *http.Request) {
	pages.Render(w, "startrack", ctrl.sessionData(w, r))
}
