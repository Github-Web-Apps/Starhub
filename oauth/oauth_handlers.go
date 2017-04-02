package oauth

import (
	"context"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// LoginHandler start the login process
func (o *Oauth) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var url = o.config.AuthCodeURL(o.state, oauth2.AccessTypeOnline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// LogoutHandler logouts current user
func (o *Oauth) LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := o.session.Get(r, o.sessionName)
		session.Values = map[interface{}]interface{}{}
		if err := session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "", http.StatusTemporaryRedirect)
	}
}

// LoginCallbackHandler deals with the login callback from github
func (o *Oauth) LoginCallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var state = r.FormValue("state")
		var code = r.FormValue("code")
		var ctx = context.Background()
		if state != o.state {
			http.Error(w, "invalid oauth state", http.StatusUnauthorized)
			return
		}
		token, err := o.config.Exchange(ctx, code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		var client = github.NewClient(o.config.Client(ctx, token))
		u, _, err := client.Users.Get(ctx, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		exists, _ := o.store.UserExist(int64(*u.ID))
		if err := o.store.SaveToken(int64(*u.ID), token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session, _ := o.session.Get(r, o.sessionName)
		session.Values["user_id"] = *u.ID
		session.Values["user_login"] = *u.Login
		session.Values["new_user"] = !exists
		if err := session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}
