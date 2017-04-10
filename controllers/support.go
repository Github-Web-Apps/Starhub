package controllers

import (
	"net/http"

	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/shared/pages"
	"github.com/gorilla/sessions"
)

// Support ctrl
type Support struct {
	Base
}

// NewSupport ctrl
func NewSupport(
	config config.Config,
	session sessions.Store,
) *Support {
	return &Support{
		Base: Base{
			config:  config,
			session: session,
		},
	}
}

// Handler handles /support
func (ctrl *Support) Handler(w http.ResponseWriter, r *http.Request) {
	pages.Render(w, "support", ctrl.sessionData(w, r))
}
