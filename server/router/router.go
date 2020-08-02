package router

import (
	"context"
	"fmt"
	"github.com/HaBaLeS/gnol/server/cache"
	"github.com/HaBaLeS/gnol/server/conversion"
	"github.com/HaBaLeS/gnol/server/dao"
	"github.com/HaBaLeS/gnol/server/session"
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type AppHandler struct {
	Router    chi.Router
	config    *util.ToolConfig
	dao       *dao.DAOHandler
	cache     *cache.ImageCache
	bgJobs    *conversion.JobRunner
	templates *template.Template
}

func NewHandler(config *util.ToolConfig, dao *dao.DAOHandler, cache *cache.ImageCache, bgj *conversion.JobRunner) *AppHandler {
	ah := &AppHandler{
		Router: chi.NewRouter(),
		config: config,
		dao:    dao,
		cache:  cache,
		bgJobs: bgj,
	}

	ah.initTemplates()
	return ah
}

func (r *AppHandler) SetupRoutes() {

	r.Router.Use(middleware.DefaultLogger)
	r.Router.Use(userSession)

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

		err = renderTemplate(tpl, w, req, cl) //TODO move template selection out!
		if err != nil {
			panic(err)
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
		err = renderTemplate(tpl, w, req, meta)
		if err != nil {
			renderError(err, w)
			return
		}

	})

	r.Router.Get("/read2/{comicId}/{imageId}", func(w http.ResponseWriter, req *http.Request) {
		comicID := chi.URLParam(req, "comicId")
		image := chi.URLParam(req, "imageId")
		num, ce := strconv.Atoi(image)
		if ce != nil {
			renderError(ce, w)
			return
		}

		//get file from cache
		var err error
		file, hit := r.cache.GetFileFromCache(comicID, num)
		if !hit {
			file, err = r.dao.GetPageImage(comicID, num)
			if err != nil {
				renderError(err, w)
				return
			}
			r.cache.AddFileToCache(file)
		}

		//as a image-provider module not the cache directly
		img, oerr := os.Open(file)
		if oerr != nil {
			renderError(oerr, w)
			return
		}

		data, rerr := ioutil.ReadAll(img)
		if rerr != nil {
			renderError(rerr, w)
			return
		}
		w.Write(data)
	})
}

func userSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("gnol")
		fmt.Printf("Session: %s\n", c)
		var us *session.UserSession
		if c != nil {
			us = session.UserSessionByID(c.Value)
		}
		if us == nil {
			fmt.Println("newSession")
			us = session.NewUserSession()
			http.SetCookie(w, &http.Cookie{Name: "gnol", Path: "/", Value: us.SessionID, Expires: time.Now().Add(time.Hour * 24)})
		}
		ctx := context.WithValue(r.Context(), "user-session", us)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserSession(ctx context.Context) *session.UserSession {
	us := ctx.Value("user-session").(*session.UserSession)
	return us
}

func renderTemplate(t *template.Template, w io.Writer, r *http.Request, pageData interface{}) error {
	us := getUserSession(r.Context())
	us.D = pageData
	return t.Execute(w, us)
}

func (r *AppHandler) initTemplates() {
	var allFiles []string
	var err error
	r.templates = template.New("root")
	r.templates = r.templates.Funcs(template.FuncMap{"mod": mod})
	if r.config.LocalResources {
		fi, _ := ioutil.ReadDir("data/template/")
		for _, file := range fi {
			filename := file.Name()
			if strings.HasSuffix(filename, ".gohtml") {
				allFiles = append(allFiles, "data/template/"+filename)
			}
		}
		r.templates, err = r.templates.ParseFiles(allFiles...)
		if err != nil {
			panic(err)
		}
	} else {
		r.templates, err = vfstemplate.ParseGlob(util.StaticAssets, r.templates, "template/*.gohtml")
	}

}

func (r *AppHandler) getTemplate(name string) (*template.Template, error) {
	r.initTemplates() //FIXME this is a DEBUG only option!!
	tpl := r.templates.Lookup(name)
	return tpl, nil
}

func renderError(e error, w http.ResponseWriter) {
	w.WriteHeader(500)
	fmt.Fprintf(w, "Error: %v", e)
	fmt.Printf("%v\n", e)
}

func mod(i, j int) bool {
	return i%j == 0
}
