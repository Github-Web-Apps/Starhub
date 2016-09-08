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
	return &Template{
		templates: template.Must(template.ParseGlob("static/*.html")),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name+".html", data)
}
