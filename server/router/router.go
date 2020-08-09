//Package router holds all controllers implementing basic Business logic for the routes
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
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//AppHandler combines Router with other submodules Implementations, like DOA, Config, Cache
type AppHandler struct {
	Router    chi.Router
	config    *util.ToolConfig
	dao       *dao.DAOHandler
	cache     *cache.ImageCache
	bgJobs    *conversion.JobRunner
	templates *template.Template
}

//NewHandler Create a new AppHandler for the Server
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

//SetupRoutes set up routing for main page and comic viewer
func (ah *AppHandler) SetupRoutes() {

	ah.Router.Use(middleware.DefaultLogger)
	ah.Router.Use(userSession)

	if ah.config.LocalResources {
		fmt.Print("Using Local resources instead of embedded\n")
		workDir, _ := os.Getwd()
		filesDir := filepath.Join(workDir, "data/")
		ah.Router.Get("/*", http.FileServer(http.Dir(filesDir)).ServeHTTP)
	} else {
		ah.Router.Get("/*", http.FileServer(util.StaticAssets).ServeHTTP)
	}

	ah.Router.Get("/comics", func(w http.ResponseWriter, req *http.Request) {
		us := getUserSession(req.Context())
		if us.ComicList == nil {
			us.ComicList = ah.dao.GetComiList(us.UserID)
		}
		ah.renderTemplate("index.gohtml", w, req, nil) //TODO move template selection out!
	})

	ah.Router.Get("/read2/{comicId}", func(w http.ResponseWriter, req *http.Request) {
		comicID := chi.URLParam(req, "comicId")
		meta, nfe := ah.dao.GetMetadata(comicID)

		if nfe != nil {
			renderError(nfe, w)
			return
		}

		ah.renderTemplate("jqviewer.gohtml", w, req, meta)
	})

	ah.Router.Get("/read2/{comicId}/{imageId}", func(w http.ResponseWriter, req *http.Request) {
		comicID := chi.URLParam(req, "comicId")
		image := chi.URLParam(req, "imageId")
		num, ce := strconv.Atoi(image)
		if ce != nil {
			renderError(ce, w)
			return
		}

		//get file from cache
		var err error
		file, hit := ah.cache.GetFileFromCache(comicID, num)
		if !hit {
			file, err = ah.dao.GetPageImage(comicID, num)
			if err != nil {
				renderError(err, w)
				return
			}
			ah.cache.AddFileToCache(file)
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
		_, re := w.Write(data)
		if re != nil {
			panic(re)
		}
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



func (ah *AppHandler) initTemplates() {
	var allFiles []string
	var err error
	ah.templates = template.New("root")
	ah.templates = ah.templates.Funcs(template.FuncMap{"mod": mod})
	if ah.config.LocalResources {
		fi, _ := ioutil.ReadDir("data/template/")
		for _, file := range fi {
			filename := file.Name()
			if strings.HasSuffix(filename, ".gohtml") {
				allFiles = append(allFiles, "data/template/"+filename)
			}
		}
		ah.templates, err = ah.templates.ParseFiles(allFiles...)
		if err != nil {
			panic(err)
		}
	} else {
		ah.templates, err = vfstemplate.ParseGlob(util.StaticAssets, ah.templates, "template/*.gohtml")
	}

}

func (ah *AppHandler) getTemplate(name string) (*template.Template, error) {
	if ah.config.LocalResources {
		//Reload templates
		ah.initTemplates()
	}
	tpl := ah.templates.Lookup(name)
	return tpl, nil
}

func (ah *AppHandler) renderTemplate(templateName string, w http.ResponseWriter, r *http.Request, pageData interface{}) {
	tpl, tlerr := ah.getTemplate(templateName)
	if tlerr != nil {
		renderError(tlerr, w)
	}
	us := getUserSession(r.Context())
	us.D = pageData
	re := tpl.Execute(w, us)
	if re != nil {
		panic(re)
	}
}

func renderError(e error, w http.ResponseWriter){
	w.WriteHeader(500)
	_, re := fmt.Fprintf(w, "Error: %v", e)
	fmt.Printf("%v\n", e)
	if re != nil {
		panic(re)
	}
}

func mod(i, j int) bool {
	return i%j == 0
}
