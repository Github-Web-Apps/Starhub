package pages

import (
	"html/template"
	"net/http"
)

// Render a given page with the given data
func Render(w http.ResponseWriter, name string, data interface{}) {
	templates, err := template.ParseFiles("static/layout.html", "static/"+name+".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := templates.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
