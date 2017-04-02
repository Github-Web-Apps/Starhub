package pages

import (
	"html/template"
	"net/http"

	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/datastore"
	"github.com/gorilla/sessions"
)

func Render(w http.ResponseWriter, name string, data interface{}) {
	templates, _ := template.ParseFiles("static/layout.html", "static/"+name+".html")
	if err := templates.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type Page struct {
	session sessions.Store
	store   datastore.Datastore
	config  config.Config
}

func New(
	config config.Config,
	store datastore.Datastore,
	session sessions.Store,
) *Page {
	return &Page{
		config:  config,
		store:   store,
		session: session,
	}
}

func (page *Page) data(r *http.Request) PageData {
	var user User
	session, _ := page.session.Get(r, page.config.SessionName)
	if !session.IsNew {
		user.ID, _ = session.Values["user_id"].(int)
		user.Login, _ = session.Values["user_login"].(string)
	}
	return PageData{
		User:     user,
		ClientID: page.config.ClientID,
	}
}

// IndexHandler handles /
func (page *Page) IndexHandler(w http.ResponseWriter, r *http.Request) {
	var data = page.data(r)
	var id = int64(data.User.ID)
	stars, err := page.store.StarCount(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	followers, err := page.store.FollowerCount(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	repos, err := page.store.RepositoryCount(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Render(w, "index", IndexData{
		PageData:        data,
		StarCount:       stars,
		FollowerCount:   followers,
		RepositoryCount: repos,
	})
}

// SupportHandler handles /support
func (page *Page) SupportHandler(w http.ResponseWriter, r *http.Request) {
	Render(w, "support", page.data(r))
}

// DonateHandler handles /donate
func (page *Page) DonateHandler(w http.ResponseWriter, r *http.Request) {
	Render(w, "donate", page.data(r))
}
