package controllers

import (
	"net/http"

	"github.com/Github-Web-Apps/Starhub/config"
	"github.com/Github-Web-Apps/Starhub/shared/pages"
	"github.com/gorilla/sessions"
)

// Apps ctrl
type Apps struct {
	Base
}

// NewApps ctrl
func NewApps(
	config config.Config,
	session sessions.Store,
) *Apps {
	return &Apps{
		Base: Base{
			config:  config,
			session: session,
		},
	}
}

// Handler handles /Apps
func (ctrl *Apps) Handler(w http.ResponseWriter, r *http.Request) {
	pages.Render(w, "apps", ctrl.sessionData(w, r))
}
