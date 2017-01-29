package oauth

import (
	"context"
	"net/http"

	"github.com/caarlos0/watchub/internal/config"
	"github.com/caarlos0/watchub/internal/datastores"
	"github.com/caarlos0/watchub/internal/dto"
	"github.com/caarlos0/watchub/internal/pages"
	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

const applicationsURL = "https://github.com/settings/connections/applications/"

// Oauth info
type Oauth struct {
	config *oauth2.Config
	store  datastores.Datastore
	state  string
}

// New oauth
func New(store datastores.Datastore, config config.Config) *Oauth {
	return &Oauth{
		config: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			Scopes:       []string{"user:email,public_repo"},
			Endpoint:     githuboauth.Endpoint,
		},
		store: store,
		state: config.OauthState,
	}
}

// Client for a given token
func (o *Oauth) Client(token *oauth2.Token) *github.Client {
	return github.NewClient(o.config.Client(context.Background(), token))
}

// Mount setup de Oauth routes
func (o *Oauth) Mount(r *mux.Router) {
	r.Methods("GET").Path("/login").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := o.config.AuthCodeURL(o.state, oauth2.AccessTypeOnline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})

	r.Methods("GET").Path("/login/callback").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var state = r.FormValue("state")
		var code = r.FormValue("code")
		if state != o.state {
			http.Error(w, "invalid oauth state", http.StatusUnauthorized)
			return
		}
		token, err := o.config.Exchange(context.Background(), code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		client := github.NewClient(o.config.Client(context.Background(), token))
		u, _, err := client.Users.Get("")
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err := o.store.SaveToken(int64(*u.ID), token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pages.Render(w, "index", dto.IndexData{
			User:                  *u.Login,
			UserID:                *u.ID,
			ChangeSubscriptionURL: applicationsURL + o.config.ClientID,
		})
	})
}
