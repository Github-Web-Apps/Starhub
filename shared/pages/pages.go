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

type Page struct {
	session sessions.Store
	config  config.Config
}

func New(config config.Config, session sessions.Store) *Page {
	return &Page{
		config:  config,
		session: session,
	}
}

func (page *Page) data(r *http.Request) dto.PageData {
	var user dto.User
	session, _ := page.session.Get(r, page.config.SessionName)
	if !session.IsNew {
		user.ID, _ = session.Values["user_id"].(int)
		user.Login, _ = session.Values["user_login"].(string)
	}
	return dto.PageData{
		User:     user,
		ClientID: page.config.ClientID,
	}
}

// IndexHandler handles /
func (page *Page) IndexHandler(w http.ResponseWriter, r *http.Request) {
	Render(w, "index", dto.IndexData{
		PageData: page.data(r),
	})
}

// SupportHandler handles /support
func (page *Page) SupportHandler(w http.ResponseWriter, r *http.Request) {
	Render(w, "support", dto.IndexData{
		PageData: page.data(r),
	})
}

// DonateHandler handles /donate
func (page *Page) DonateHandler(w http.ResponseWriter, r *http.Request) {
	Render(w, "donate", dto.IndexData{
		PageData: page.data(r),
	})
}
