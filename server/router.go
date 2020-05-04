package server

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type AppHandler struct {
	router chi.Router
	session *Session
}
func NewHandler(s *Session) (*AppHandler) {
	return &AppHandler{
		session: s,
		router: chi.NewRouter(),
	}
}

func (r *AppHandler) SetupRoutes(){

	r.router.Use(middleware.DefaultLogger)
	r.router.Get("/echo/*", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Echo: %s", r.URL.Path)
	})

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "static/")
	r.router.Get("/*", http.FileServer(http.Dir(filesDir)).ServeHTTP)


	r.router.Get("/comics", func(w http.ResponseWriter, r *http.Request) {
		cl := GetComiList()
		t := template.New("index.gohtml")
		t.Funcs(template.FuncMap{"mod": mod})
		tpl, err := t.ParseFiles("template/index.gohtml")
		if err != nil {
			panic(err)
		}

		err = tpl.Execute(w, cl)
		if err != nil {
			panic(err)
		}
	})

}

	func mod(i, j int) bool {
		return i%j == 0
	}