package auth

import (
	"celme/blank"
	helper "celme/helpers"
	"net/http"
	"path"
)

// LoginHandler
// возвращает обработчик запросов отображения формы авторизации
func LoginHandler(mongoURI, dbName string) func(http.ResponseWriter, *http.Request) {
	base := "/blank/views/login"
	templates := []string{
		path.Join("/blank/views", "pure"),
		"content",
	}
	return func(w http.ResponseWriter, r *http.Request) {
		getTmpl := blank.TemplateLoader(base, helper.Assets)
		tmpl := getTmpl(templates...)
		data := blank.NewRenderData()
		tmpl.ExecuteTemplate(w, "pure", data)
	}
}
