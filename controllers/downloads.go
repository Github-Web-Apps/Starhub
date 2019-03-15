package controllers

import (
	"net/http"

	"github.com/Github-Web-Apps/Starhub/config"
	"github.com/Github-Web-Apps/Starhub/shared/pages"
	"github.com/gorilla/sessions"
)

// Downloads ctrl
type Downloads struct {
	Base
}

// NewDownloads ctrl
func NewDownloads(
	config config.Config,
	session sessions.Store,
) *Downloads {
	return &Downloads{
		Base: Base{
			config:  config,
			session: session,
		},
	}
}

// Handler handles /Downloads
func (ctrl *Downloads) Handler(w http.ResponseWriter, r *http.Request) {
	pages.Render(w, "downloads", ctrl.sessionData(w, r))
}
