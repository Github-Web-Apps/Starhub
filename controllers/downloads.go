package controllers

import (
	"net/http"

	"github.com/Intika-Web-Apps/Starhub-Notifier/config"
	"github.com/Intika-Web-Apps/Starhub-Notifier/shared/pages"
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
