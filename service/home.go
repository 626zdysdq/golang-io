package service

import (
	"net/http"
	"text/template"

	"github.com/unrolled/render"
)

//请求主页之后，发送主页面
func home(w http.ResponseWriter, r *http.Request) {

	t := template.Must(template.ParseFiles("templates/login.html"))
	t.Execute(w, nil)
}

func apiTestHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct {
			WELCOME string `json:"WELCOME"`
		}{WELCOME: "Hello"})
	}
}
