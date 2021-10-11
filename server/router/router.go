//Package router holds all controllers implementing basic Business logic for the routes
package router

import (
	"context"
	"fmt"
	"github.com/HaBaLeS/gnol/server/cache"
	"github.com/HaBaLeS/gnol/server/jobs"
	"github.com/HaBaLeS/gnol/server/session"
	"github.com/HaBaLeS/gnol/server/storage"
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//AppHandler combines Router with other submodules Implementations, like DOA, Config, Cache
type AppHandler struct {
	Router    chi.Router
	config    *util.ToolConfig
	dao		*storage.DAO
	cache     *cache.ImageCache
	bgJobs    *jobs.JobRunner
	templates *template.Template
}

//NewHandler Create a new AppHandler for the Server
func NewHandler(config *util.ToolConfig, cache *cache.ImageCache, bgj *jobs.JobRunner, dao	*storage.DAO) *AppHandler {
	ah := &AppHandler{
		Router: chi.NewRouter(),
		config: config,
		cache:  cache,
		bgJobs: bgj,
		dao: dao,
	}

	ah.initTemplates()
	return ah
}

//Routes defines all routes for /user and below.
//this path cares about UserManagement
func (ah *AppHandler) Routes() {

	//Define global middleware
	ah.Router.Use(middleware.DefaultLogger)
	ah.Router.Use(userSession)

	ah.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/comics", 301)
	})



	//Handle static Resources
	if ah.config.LocalResources {
		fmt.Print("Using Local resources instead of embedded\n")
		workDir, _ := os.Getwd()
		filesDir := filepath.Join(workDir, "data/")
		ah.Router.Get("/*", http.FileServer(http.Dir(filesDir)).ServeHTTP)
	} else {
		ah.Router.Get("/*", http.FileServer(util.StaticAssets).ServeHTTP)
	}

	//Define users
	ah.Router.Route("/users", func(r chi.Router) { //FIXME remove s in users
		r.Get("/", ah.listUsers())
		r.Post("/", ah.createUser())
		r.Route("/{userID}", func(r chi.Router) {
			r.Get("/", ah.getUser())
			r.Put("/", ah.updateUser())
			r.Delete("/", ah.deleteUser())
		})
		r.Get("/create", ah.serveTemplate("create_user.gohtml", nil))
		r.Get("/login", ah.serveTemplate("login_user.gohtml", nil))
		r.Post("/login", ah.loginUser())
		r.Get("/logout", ah.logoutUser())
	})

	//Define Uploads
	ah.Router.Route("/upload", func(r chi.Router) {
		r.Get("/archive", ah.serveTemplate("upload_archive.gohtml",nil))
		r.Get("/pdf", ah.serveTemplate("upload_pdf.gohtml",nil))
		r.Get("/url", ah.serveTemplate("upload_url.gohtml",nil))
		r.Post("/archive", ah.uploadArchive())
		r.Post("/url", ah.uploadUrl())
		r.Post("/pdf", ah.uploadPdf())
	})

	//Define Comic
	ah.Router.Route("/comics", func(r chi.Router) {
		r.Get("/", ah.comicsList())
		r.Get("/{comicId}", ah.comicsLoad())
		r.Get("/{comicId}/{imageId}", ah.comicsPageImage())
	})

	//Define Series
	ah.Router.Route("/series", func(r chi.Router) {
		r.Get("/", ah.seriesList())
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
