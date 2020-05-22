package server

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"playground.dahoam/util"
)

type AppHandler struct {
	router  chi.Router
	session *Session
}

func NewHandler(s *Session) *AppHandler {
	return &AppHandler{
		session: s,
		router:  chi.NewRouter(),
	}
}

func (r *AppHandler) SetupRoutes() {

	r.router.Use(middleware.DefaultLogger)
	r.router.Get("/echo/*", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Echo: %s", r.URL.Path)
	})

	if r.session.config.LocalResources {
		fmt.Print("Using Local resources instead of embedded\n")
		workDir, _ := os.Getwd()
		filesDir := filepath.Join(workDir, "static/")
		r.router.Get("/*", http.FileServer(http.Dir(filesDir)).ServeHTTP)
	} else {
		r.router.Get("/*", http.FileServer(util.StaticAssets).ServeHTTP)
	}

	r.router.Get("/comics", func(w http.ResponseWriter, req *http.Request) {
		cl := r.session.dao.GetComiList()

		tpl, err := r.getTemplate("index.gohtml")
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

		tpl, err := r.getTemplate("view2.gohtml")
		err = tpl.Execute(w, meta)

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

		tpl, err := r.getTemplate("jqviewer.gohtml")
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
			renderError(err, w)
			return
		}

		w.Write(data)
	})

	r.router.Get("/read2/{comicId}/{imageId}", func(w http.ResponseWriter, req *http.Request) {
		comicId := chi.URLParam(req, "comicId")
		image := chi.URLParam(req, "imageId")
		data, err := r.session.dao.getPageImage(comicId, image)
		if err != nil {
			renderError(err, w)
			return
		}

		w.Write(data)
	})
}

func (r *AppHandler) getTemplate(name string) (*template.Template, error) {
	var tpl *template.Template
	var err error
	t := template.New(name)
	t.Funcs(template.FuncMap{"mod": mod})
	if r.session.config.LocalResources {
		tpl, err = t.ParseFiles("data/template/" + name)
	} else {
		tpl, err = vfstemplate.ParseFiles(util.StaticAssets, t, "template/"+name)
	}

	if err != nil {
		return nil, err
	}
	return tpl, nil
}

func renderError(e error, w http.ResponseWriter) {
	fmt.Fprintf(w, "Error: %v", e)
}

func mod(i, j int) bool {
	return i%j == 0
}
