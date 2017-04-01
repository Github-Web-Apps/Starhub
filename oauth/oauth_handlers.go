package oauth

import (
	"context"
	"net/http"

	"github.com/caarlos0/watchub/shared/dto"
	"github.com/caarlos0/watchub/shared/pages"
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
		if err := o.store.SaveToken(int64(*u.ID), token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session, _ := o.session.Get(r, o.sessionName)
		session.Values["user_id"] = *u.ID
		session.Values["user_login"] = *u.Login
		if err := session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pages.Render(w, "index", dto.IndexData{
			User: dto.User{
				ID:    *u.ID,
				Login: *u.Login,
			},
			ClientID: o.config.ClientID,
		})
	}
}
