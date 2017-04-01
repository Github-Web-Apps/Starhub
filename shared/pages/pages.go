package pages

import (
	"html/template"
	"net/http"

	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/shared/dto"
	"github.com/gorilla/sessions"
)

func Render(w http.ResponseWriter, name string, data interface{}) {
	templates, _ := template.ParseFiles("static/layout.html", "static/"+name+".html")
	if err := templates.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type GenericPage struct {
	session sessions.Store
	config  config.Config
	name    string
}

func New(
	config config.Config,
	session sessions.Store,
	page string,
) *GenericPage {
	return &GenericPage{
		config:  config,
		session: session,
		name:    page,
	}
}

func (page *GenericPage) Handler(w http.ResponseWriter, r *http.Request) {
	session, _ := page.session.Get(r, page.config.SessionName)
	var user dto.User
	if !session.IsNew {
		user.ID, _ = session.Values["user_id"].(int)
		user.Login, _ = session.Values["user_login"].(string)
	}
	Render(w, page.name, dto.IndexData{
		User:     user,
		ClientID: page.config.ClientID,
	})
}
