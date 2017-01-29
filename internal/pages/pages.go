package pages

import (
	"html/template"
	"io"
)

func Render(w io.Writer, name string, data interface{}) error {
	templates, _ := template.ParseFiles("static/layout.html", "static/"+name+".html")
	return templates.ExecuteTemplate(w, "layout", data)
}
