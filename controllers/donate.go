package controllers

import (
	"net/http"

	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/shared/pages"
	"github.com/gorilla/sessions"
)

// Donate ctrl
type Donate struct {
	Base
}

// NewDonate ctrl
func NewDonate(
	config config.Config,
	session sessions.Store,
) *Donate {
	return &Donate{
		Base: Base{
			config:  config,
			session: session,
		},
	}
}

// Handler handles /donate
func (ctrl *Donate) Handler(w http.ResponseWriter, r *http.Request) {
	pages.Render(w, "donate", ctrl.sessionData(w, r))
}
