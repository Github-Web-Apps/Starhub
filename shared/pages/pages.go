package pages

import (
	"html/template"
	"net/http"
)

func Render(w http.ResponseWriter, name string, data interface{}) {
	templates, _ := template.ParseFiles("static/layout.html", "static/"+name+".html")
	if err := templates.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
