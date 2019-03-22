package controllers

import (
	"net/http"

	"github.com/Github-Web-Apps/Starhub/config"
	"github.com/Github-Web-Apps/Starhub/shared/pages"
	"github.com/gorilla/sessions"
)

// Profile ctrl
type Profile struct {
	Base
}

// NewProfile ctrl
func NewProfile(
	config config.Config,
	session sessions.Store,
) *Profile {
	return &Profile{
		Base: Base{
			config:  config,
			session: session,
		},
	}
}

// Handler handles /Profile
func (ctrl *Profile) Handler(w http.ResponseWriter, r *http.Request) {
	pages.Render(w, "profile", ctrl.sessionData(w, r))
}
