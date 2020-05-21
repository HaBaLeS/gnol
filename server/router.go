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


	r.router.Get("/comics", func(w http.ResponseWriter, req *http.Request) {
		cl := r.session.dao.GetComiList()
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

	r.router.Get("/read/{comicId}", func(w http.ResponseWriter, req *http.Request) {
		comicId := chi.URLParam(req, "comicId")
		meta, nfe := r.session.dao.getMetadata(comicId)

		if nfe != nil {
			renderError(nfe, w)
			return
		}

		t := template.New("view2.gohtml")
		tpl, err := t.ParseFiles("template/view2.gohtml")
		if err != nil {
			renderError(err, w)
			return
		}

		err = tpl.Execute(w, meta)
		if err != nil {
			renderError(err, w)
			return
		}
	})

	r.router.Get("/read2/{comicId}", func(w http.ResponseWriter, req *http.Request) {
		comicId := chi.URLParam(req, "comicId")
		meta, nfe := r.session.dao.getMetadata(comicId)

		if nfe != nil {
			renderError(nfe, w)
			return
		}

		t := template.New("jqviewer.gohtml")
		tpl, err := t.ParseFiles("template/jqviewer.gohtml")
		if err != nil {
			renderError(err, w)
			return
		}

		err = tpl.Execute(w, meta)
		if err != nil {
			renderError(err, w)
			return
		}

	})
	
	r.router.Get("/read/{comicId}/{imageId}", func(w http.ResponseWriter, req *http.Request) {
		comicId := chi.URLParam(req, "comicId")
		image := chi.URLParam(req, "imageId")
		data, err := r.session.dao.getPageImage(comicId, image)
		if err != nil {
			renderError(err,w)
			return
		}

		w.Write(data)
	})

	r.router.Get("/read2/{comicId}/{imageId}", func(w http.ResponseWriter, req *http.Request) {
		comicId := chi.URLParam(req, "comicId")
		image := chi.URLParam(req, "imageId")
		data, err := r.session.dao.getPageImage(comicId, image)
		if err != nil {
			renderError(err,w)
			return
		}

		w.Write(data)
	})
}

func renderError(e error, w http.ResponseWriter) {
	fmt.Fprintf(w, "Error: %v", e)
}

func mod(i, j int) bool {
	return i%j == 0
}