package controllers

import (
	"net/http"

	"github.com/Github-Web-Apps/Starhub/config"
	"github.com/Github-Web-Apps/Starhub/shared/pages"
	"github.com/gorilla/sessions"
)

// Tools ctrl
type Tools struct {
	Base
}

// NewTools ctrl
func NewTools(
	config config.Config,
	session sessions.Store,
) *Tools {
	return &Tools{
		Base: Base{
			config:  config,
			session: session,
		},
	}
}

// Handler handles /Downloads
func (ctrl *Tools) Handler(w http.ResponseWriter, r *http.Request) {
	pages.Render(w, "tools", ctrl.sessionData(w, r))
}
