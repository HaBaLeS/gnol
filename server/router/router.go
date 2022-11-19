// Package router holds all controllers implementing basic Business logic for the routes
package router

import (
	"encoding/gob"
	"fmt"
	"github.com/HaBaLeS/gnol/data/static"
	template2 "github.com/HaBaLeS/gnol/data/template"
	"github.com/HaBaLeS/gnol/server/cache"
	"github.com/HaBaLeS/gnol/server/gnolsession"
	"github.com/HaBaLeS/gnol/server/jobs"
	"github.com/HaBaLeS/gnol/server/storage"
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// AppHandler combines Router with other submodules Implementations, like DOA, Config, Cache
type AppHandler struct {
	Router    *gin.Engine
	config    *util.ToolConfig
	dao       *storage.DAO
	cache     *cache.ImageCache
	bgJobs    *jobs.JobRunner
	templates *template.Template
	web       *webauthn.WebAuthn
}

// NewHandler Create a new AppHandler for the Server
func NewHandler(config *util.ToolConfig, cache *cache.ImageCache, bgj *jobs.JobRunner, dao *storage.DAO) *AppHandler {
	ah := &AppHandler{
		Router: gin.Default(), //Fixme, don't use defaults
		config: config,
		cache:  cache,
		bgJobs: bgj,
		dao:    dao,
	}

	web, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "GNOL Online",               // Display Name for your site
		RPID:          config.WebAuthnHostname,     // Generally the FQDN for your site
		RPOrigin:      config.WebAuthnOriginURL,    // The origin URL for WebAuthn requests
		RPIcon:        "http://localhost/logo.png", // Optional icon URL for your site
	})
	if err != nil {
		fmt.Println(err)
	}
	ah.web = web

	ah.initTemplates()
	return ah
}

// Routes defines all routes for /user and below.
// this path cares about UserManagement
func (ah *AppHandler) Routes() {

	gob.Register(gnolsession.UserSession{})

	//Define global middleware
	store := gnolsession.NewGnolSessionStore()
	ah.Router.Use(sessions.Sessions("gnolsession", store))
	ah.Router.Use(userSessionMiddleware)

	ah.Router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/series")
	})

	//Handle static Resources
	if ah.config.LocalResources {
		fmt.Print("Using Local resources instead of embedded\n")
		workDir, _ := os.Getwd()
		filesDir := filepath.Join(workDir, "data/static/")
		ah.Router.StaticFS("/static", http.Dir(filesDir))
	} else {
		ah.Router.StaticFS("/static", http.FS(static.Embedded))

	}

	//Define users
	user := ah.Router.Group("/users")
	{
		user.GET("/", ah.listUsers)
		user.POST("/", ah.createUser)
		user.GET("/create", ah.serveTemplate("register.gohtml", nil))
		user.GET("/login", ah.serveTemplate("login_user.gohtml", nil))
		user.POST("/login", ah.loginUser)
		user.GET("/logout", ah.logoutUser)
	}

	stng := ah.Router.Group("/setting")
	{
		stng.Use(requireAuth)
		stng.GET("/api-token", ah.APIToken)
	}

	wn := ah.Router.Group("/webauthn")
	{
		wn.GET("/", ah.webAuthnIndex)
		wn.GET("/:userID", ah.BeginRegistration)
		wn.POST("/add", ah.FinishRegistration)
		wn.GET("/assertion/:userID", ah.BeginLogin)
		wn.POST("/assertion", ah.FinishLogin)
	}

	//Define Uploads
	up := ah.Router.Group("/upload")
	{
		up.Use(requireAuth)
		up.GET("/archive", ah.serveTemplate("upload_archive.gohtml", nil))
		up.GET("/pdf", ah.serveTemplate("upload_pdf.gohtml", nil))
		up.GET("/url", ah.serveTemplate("upload_url.gohtml", nil))
		up.POST("/archive", ah.uploadArchive)
		up.POST("/url", ah.uploadUrl)
		up.POST("/pdf", ah.uploadPdf)
	}

	//Define Comic
	cm := ah.Router.Group("/comics")
	{
		cm.Use(requireAuth)
		cm.GET("/", ah.comicsList)
		cm.GET("/:comicId", ah.comicsLoad)
		cm.GET("/:comicId/edit", ah.comicsEdit)
		cm.POST("/:comicId/edit", ah.updateComic)
		cm.GET("/:comicId/continue/:lastpage", ah.comicsLoad)
		cm.GET("/:comicId/:imageId", ah.comicsPageImage)
		cm.PUT("/last/:comicId/:lastpage", ah.comicSetLastPage)
		cm.DELETE("/delete/:comicId", ah.deleteComic)
	}

	//Define Series
	srs := ah.Router.Group("/series")
	{
		srs.Use(requireAuth)
		srs.GET("/", ah.seriesList)
		srs.GET("/:seriesID", ah.comicsInSeriesList)
		srs.GET("/create", ah.serveTemplate("series_create.gohtml", nil))
		srs.POST("/create", ah.createSeries)
	}

	api := ah.Router.Group("/api")
	{
		api.Use(ah.requireAPIToken)
		api.GET("/list", ah.apiListComics)
		api.GET("/series", ah.apiSeries)
		api.POST("/upload", ah.apiUploadComic)
		/*
			get /api/series/list
			post /api/series/create
			post /api/series/update
		*/
	}

	ah.Router.NoRoute(func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/favicon.ico" {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Redirect(http.StatusTemporaryRedirect, "/series")
	})
}

func requireAuth(ctx *gin.Context) {
	ssn := sessions.Default(ctx)
	gnoluser := ssn.Get("user-session")
	if gnoluser != nil && gnoluser.(*gnolsession.UserSession).IsLoggedIn() {
		ctx.Next()
		return
	}
	ctx.Redirect(http.StatusTemporaryRedirect, "/users/login")
	ctx.Abort()
}

func userSessionMiddleware(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var us *gnolsession.UserSession
	if session.Get("user-session") == nil {
		fmt.Println("newSession")
		us = gnolsession.NewUserSession()
		session.Set("user-session", us)
	}
	ctx.Next()

	if e := session.Save(); e != nil {
		panic(e)
	}

}

func getUserSession(ctx *gin.Context) *gnolsession.UserSession {
	s := sessions.Default(ctx)
	us := s.Get("user-session").(*gnolsession.UserSession)
	return us
}

func updateUSerSession(ctx *gin.Context, us *gnolsession.UserSession) {
	s := sessions.Default(ctx)
	s.Set("user-session", us)
	s.Save()
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
		ah.templates, err = ah.templates.ParseFS(template2.Embedded, "*.gohtml")
		if err != nil {
			panic(err)
		}
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

func (ah *AppHandler) renderTemplate(templateName string, ctx *gin.Context, pageData interface{}) {
	tpl, tlerr := ah.getTemplate(templateName)
	if tlerr != nil {
		renderError(tlerr, ctx.Writer)
	}
	us := getUserSession(ctx)
	us.D = pageData
	re := tpl.Execute(ctx.Writer, us)
	if re != nil {
		panic(re)
	}
}

func renderError(e error, w http.ResponseWriter) {
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
