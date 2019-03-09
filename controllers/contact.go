package controllers

import (
	"net/http"

	"github.com/Intika-Web-Apps/Watchub-Mirror/config"
	"github.com/Intika-Web-Apps/Watchub-Mirror/shared/pages"
	"github.com/gorilla/sessions"
)

// Contact ctrl
type Contact struct {
	Base
}

// NewContact ctrl
func NewContact(
	config config.Config,
	session sessions.Store,
) *Contact {
	return &Contact{
		Base: Base{
			config:  config,
			session: session,
		},
	}
}

// Handler handles /Contact
func (ctrl *Contact) Handler(w http.ResponseWriter, r *http.Request) {
	pages.Render(w, "contact", ctrl.sessionData(w, r))
}
