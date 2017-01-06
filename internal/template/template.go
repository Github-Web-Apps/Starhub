package template

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

func New(glob string) *Template {
	return &Template{}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	templates, _ := template.ParseFiles("static/layout.html", "static/"+name+".html")
	return templates.ExecuteTemplate(w, "layout", data)
}
