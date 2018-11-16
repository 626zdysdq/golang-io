package service

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {
	formatter := render.New()

	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, formatter)

	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	webRoot := os.Getenv("WEBROOT")
	//fmt.Print(webRoot + "a")
	if len(webRoot) == 0 {
		if root, err := os.Getwd(); err != nil {
			panic("Could not retrive working directory")
		} else {
			webRoot = root
			//`fmt.Println(root)
		}
	}
	//实现对静态文件和js的引用
	mx.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(webRoot+"/add/"))))
	//表单提交与访问
	mx.HandleFunc("/", home).Methods("GET")
	mx.HandleFunc("/", login).Methods("POST")
	//unknown的错误显示
	mx.HandleFunc("/unknown", NotImplemented).Methods("GET")
	mx.HandleFunc("/api/test", apiTestHandler(formatter)).Methods("GET")

}
