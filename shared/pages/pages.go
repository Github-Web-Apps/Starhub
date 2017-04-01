package pages

import (
	"html/template"
	"net/http"

	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/shared/dto"
)

func Render(w http.ResponseWriter, name string, data interface{}) {
	templates, _ := template.ParseFiles("static/layout.html", "static/"+name+".html")
	if err := templates.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type GenericPage struct {
	config config.Config
	name   string
}

func New(config config.Config, name string) *GenericPage {
	return &GenericPage{
		config: config,
		name:   name,
	}
}

func (gp *GenericPage) Handler(w http.ResponseWriter, r *http.Request) {
	Render(w, gp.name, dto.IndexData{
		User: dto.User{
			Login: "moises",
		},
		ClientID: gp.config.ClientID,
	})
}
