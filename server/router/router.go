package router

import (
	"fmt"
	"github.com/HaBaLeS/gnol/server/dao"
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type AppHandler struct {
	Router chi.Router
	config *util.ToolConfig
	dao    *dao.DAOHandler
}

func NewHandler(config *util.ToolConfig, dao *dao.DAOHandler) *AppHandler {
	return &AppHandler{
		Router: chi.NewRouter(),
		config: config,
		dao:    dao,
	}
}

func (r *AppHandler) SetupRoutes() {

	r.Router.Use(middleware.DefaultLogger)
	r.Router.Get("/echo/*", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Echo: %s", r.URL.Path)
	})

	if r.config.LocalResources {
		fmt.Print("Using Local resources instead of embedded\n")
		workDir, _ := os.Getwd()
		filesDir := filepath.Join(workDir, "data/")
		r.Router.Get("/*", http.FileServer(http.Dir(filesDir)).ServeHTTP)
	} else {
		r.Router.Get("/*", http.FileServer(util.StaticAssets).ServeHTTP)
	}

	r.Router.Get("/comics", func(w http.ResponseWriter, req *http.Request) {
		cl := r.dao.GetComiList()

		tpl, err := r.getTemplate("index.gohtml")
		if err != nil {
			panic(err)
		}

		err = tpl.Execute(w, cl)
		if err != nil {
			panic(err)
		}
	})

	r.Router.Get("/read/{comicId}", func(w http.ResponseWriter, req *http.Request) {
		comicId := chi.URLParam(req, "comicId")
		meta, nfe := r.dao.GetMetadata(comicId)

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

	r.Router.Get("/read2/{comicId}", func(w http.ResponseWriter, req *http.Request) {
		comicId := chi.URLParam(req, "comicId")
		meta, nfe := r.dao.GetMetadata(comicId)

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

	r.Router.Get("/read/{comicId}/{imageId}", func(w http.ResponseWriter, req *http.Request) {
		comicId := chi.URLParam(req, "comicId")
		image := chi.URLParam(req, "imageId")
		data, err := r.dao.GetPageImage(comicId, image)
		if err != nil {
			renderError(err, w)
			return
		}

		w.Write(data)
	})

	r.Router.Get("/read2/{comicId}/{imageId}", func(w http.ResponseWriter, req *http.Request) {
		/*comicId := chi.URLParam(req, "comicId")
		image := chi.URLParam(req, "imageId")
		num, ce := strconv.Atoi(image)
		if ce != nil {
			renderError(ce, w)
			return
		}

		//as a image-provider module not the cache directly
		r.session.cache.AddFileToCache("")
		loader, err := //GetImage(comicId, num)
		data, err2 := loader()
		if err2 != nil {
			renderError(err, w)
			return
		}

		w.Write(data)*/
	})
}

func (r *AppHandler) getTemplate(name string) (*template.Template, error) {
	var tpl *template.Template
	var err error
	t := template.New(name)
	t.Funcs(template.FuncMap{"mod": mod})
	if r.config.LocalResources {
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
	w.WriteHeader(500)
	fmt.Fprintf(w, "Error: %v", e)
}

func mod(i, j int) bool {
	return i%j == 0
}