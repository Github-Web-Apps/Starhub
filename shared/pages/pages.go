package pages

import (
	"html/template"
	"net/http"

	"github.com/apex/log"
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

func (page *Page) data(w http.ResponseWriter, r *http.Request) PageData {
	var user User
	session, _ := page.session.Get(r, page.config.SessionName)
	if !session.IsNew {
		user.ID, _ = session.Values["user_id"].(int)
		user.Login, _ = session.Values["user_login"].(string)
		user.IsNew, _ = session.Values["new_user"].(bool)
		delete(session.Values, "new_user")
		if err := session.Save(r, w); err != nil {
			log.WithError(err).Error("failed to update session")
		}
	}
	return PageData{
		User:     user,
		ClientID: page.config.ClientID,
	}
}

// IndexHandler handles /
func (page *Page) IndexHandler(w http.ResponseWriter, r *http.Request) {
	var data = IndexData{
		PageData: page.data(w, r),
	}
	if data.User.ID > 0 {
		var err error
		var id = int64(data.User.ID)
		data.StarCount, err = page.store.StarCount(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.FollowerCount, err = page.store.FollowerCount(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.RepositoryCount, err = page.store.RepositoryCount(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	Render(w, "index", data)
}

// SupportHandler handles /support
func (page *Page) SupportHandler(w http.ResponseWriter, r *http.Request) {
	Render(w, "support", page.data(w, r))
}

// DonateHandler handles /donate
func (page *Page) DonateHandler(w http.ResponseWriter, r *http.Request) {
	Render(w, "donate", page.data(w, r))
}
