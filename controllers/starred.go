package controllers

import (
	"net/http"

	"github.com/Github-Web-Apps/Starhub/config"
	"github.com/Github-Web-Apps/Starhub/shared/pages"
	"github.com/gorilla/sessions"
)

// Starred ctrl
type Starred struct {
	Base
}

// NewStarred ctrl
func NewStarred(
	config config.Config,
	session sessions.Store,
) *Starred {
	return &Starred{
		Base: Base{
			config:  config,
			session: session,
		},
	}
}

// Handler handles /Downloads
func (ctrl *Starred) Handler(w http.ResponseWriter, r *http.Request) {
	pages.Render(w, "starred", ctrl.sessionData(w, r))
}
