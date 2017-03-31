package oauth

import (
	"context"
	"net/http"

	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/datastore"
	"github.com/caarlos0/watchub/shared/dto"
	"github.com/caarlos0/watchub/shared/pages"
	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

// Oauth info
type Oauth struct {
	config *oauth2.Config
	store  datastore.Datastore
	state  string
}

// New oauth
func New(store datastore.Datastore, config config.Config) *Oauth {
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

// ClientFrom for a given string token
func (o *Oauth) ClientFrom(ctx context.Context, tokenStr string) *github.Client {
	return o.Client(ctx, token.FromJson(tokenStr))
}

// Client for a given token
// TODO check if this will still be used
func (o *Oauth) Client(ctx context.Context, token *oauth2.Token) *github.Client {
	return github.NewClient(o.config.Client(ctx, token))
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
		ctx := context.Background()
		token, err := o.config.Exchange(ctx, code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		client := github.NewClient(o.config.Client(ctx, token))
		u, _, err := client.Users.Get(ctx, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err := o.store.SaveToken(int64(*u.ID), token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pages.Render(w, "index", dto.IndexData{
			User:     *u.Login,
			UserID:   *u.ID,
			ClientID: o.config.ClientID,
		})
	})
}
