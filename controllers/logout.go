package controllers

import (
	"net/http"

	"github.com/Intika-Web-Apps/Watchub-Mirror/config"
	"github.com/gorilla/sessions"
)

// Logout ctrl
type Logout struct {
	Base
}

// NewLogout ctrl
func NewLogout(
	config config.Config,
	session sessions.Store,
) *Logout {
	return &Logout{
		Base: Base{
			config:  config,
			session: session,
		},
	}
}

// Handler handles /
func (ctrl *Logout) Handler(w http.ResponseWriter, r *http.Request) {
	session, _ := ctrl.session.Get(r, ctrl.config.SessionName)
	session.Values = map[interface{}]interface{}{}
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "", http.StatusTemporaryRedirect)
}
