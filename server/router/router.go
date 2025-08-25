// Package router holds all controllers implementing basic Business logic for the routes
package router

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/HaBaLeS/gnol/data/static"
	template2 "github.com/HaBaLeS/gnol/data/template"
	"github.com/HaBaLeS/gnol/docs"
	"github.com/HaBaLeS/gnol/server/cache"
	"github.com/HaBaLeS/gnol/server/jobs"
	"github.com/HaBaLeS/gnol/server/storage"
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type GnolContext struct {
	Issue      *storage.Comic
	Series     *storage.Series
	ComicList  []*storage.Comic
	SeriesList []*storage.Series
	UserList   []*storage.User
	Session    *storage.GnolSession
	Flash      string
	UserInfo   *storage.User
}

func NewGnolContext(gs *storage.GnolSession) *GnolContext {
	return &GnolContext{
		Session: gs,
	}
}

// AppHandler combines Router with other submodules Implementations, like DOA, Config, Cache
type AppHandler struct {
	Router    *gin.Engine
	config    *util.ToolConfig
	dao       *storage.DAO
	cache     *cache.ImageCache
	bgJobs    *jobs.JobRunner
	templates *template.Template
	//web       *webauthn.WebAuthn
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

	/*web, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "GNOL Online",               // Display Name for your site
		RPID:          config.WebAuthnHostname,     // Generally the FQDN for your site
		RPOrigin:      config.WebAuthnOriginURL,    // The origin URL for WebAuthn requests
		RPIcon:        "http://localhost/logo.png", // Optional icon URL for your site
	})
	if err != nil {
		fmt.Println(err)
	}
	ah.web = web*/

	ah.initTemplates()
	return ah
}

// Routes defines all routes for /user and below.
// this path cares about UserManagement
func (ah *AppHandler) Routes() {

	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = "localhost:8666"
	ah.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Define global middleware
	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		MaxAge:   86400 * 30,
		Path:     "/",
		Secure:   false, //fixme only for dev mode!!
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	ah.Router.Use(sessions.Sessions("gnol-session-id", store))
	ah.Router.Use(userSessionMiddleware)

	ah.Router.GET("/", redirect("/series"))

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
		user.GET("/create", ah.serveTemplate("register.gohtml"))
		user.GET("/login", ah.serveTemplate("login_user.gohtml"))

		user.POST("/login", ah.loginUser)

		authUser := user.Use(ah.requireAuth)
		authUser.GET("/logout", ah.logoutUser)
		authUser.GET("/profile", ah.serveTemplate("user_profile.gohtml")) //FIXME add support for template folder structures. Currently only flat folder is supported in  initTemplates()

	}

	stng := ah.Router.Group("/setting")
	{
		stng.Use(ah.requireAuth)
		stng.GET("/api-token", ah.APIToken)
	}

	/*wn := ah.Router.Group("/webauthn")
	{
		wn.GET("/", ah.webAuthnIndex)
		wn.GET("/:userID", ah.BeginRegistration)
		wn.POST("/add", ah.FinishRegistration)
		wn.GET("/assertion/:userID", ah.BeginLogin)
		wn.POST("/assertion", ah.FinishLogin)
	}*/

	//Define Uploads
	up := ah.Router.Group("/upload")
	{
		up.Use(ah.requireAuth)

		up.GET("/archive", func(context *gin.Context) {
			panic("Not implementeded spectial Render Func!")
			//sl := ah.dao.AllSeries()
			//ah.renderTemplate("upload_archive.gohtml", context, sl)
		})
		up.GET("/pdf", ah.serveTemplate("upload_pdf.gohtml"))
		up.GET("/url", ah.serveTemplate("upload_url.gohtml"))
		up.POST("/archive", ah.uploadArchive)
		up.POST("/url", ah.uploadUrl)
		up.POST("/pdf", ah.uploadPdf)
	}

	//Define Comic
	cm := ah.Router.Group("/comics")
	{
		cm.Use(ah.requireAuth)
		cm.GET("/", redirect("/series"))
		cm.GET("/:comicId", ah.comicsLoad)
		cm.GET("/:comicId/edit", ah.comicsEdit)
		cm.POST("/:comicId/edit", ah.updateComic)
		cm.GET("/:comicId/continue/:lastpage", ah.comicsLoad)
		cm.GET("/:comicId/:imageId", ah.comicsPageImage)
		cm.PUT("/last/:comicId/:lastpage", ah.comicSetLastPage)
		cm.DELETE("/remove/:comicId", ah.removeComic)
		cm.DELETE("/delete/:comicId", ah.deleteComic)
	}

	//Define Comic
	sh := ah.Router.Group("/share")
	{
		sh.Use(ah.requireAuth)
		sh.PUT("/comic/:comicId/:userId", ah.shareComic)
		sh.PUT("/series/:seriesId/:userId", ah.shareSeries)
	}

	//Define Series
	srs := ah.Router.Group("/series")
	{
		srs.Use(ah.requireAuth)
		srs.GET("/", ah.seriesList)
		srs.GET("/:seriesId", ah.comicsInSeriesList)
		srs.GET("/create", ah.serveTemplate("series_create.gohtml"))
		srs.POST("/create", ah.createSeries)
		srs.GET("/:seriesId/edit", ah.seriesEdit)    //Render Edit Page
		srs.POST("/:seriesId/edit", ah.updateSeries) //FIXME this should share stuff witl API!!
	}

	//Define Series
	srsNsfw := ah.Router.Group("/seriesNSFW")
	{
		srsNsfw.Use(ah.requireAuth)
		srsNsfw.GET("/", ah.seriesListNSFW)
	}

	api := ah.Router.Group("/api")
	{
		api.Use(ah.requireAPIToken)
		api.GET("/list", ah.apiListComics)
		api.GET("/series", ah.apiSeries)
		api.GET("/checkhash/:hash", ah.apiCheckHash)
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

func redirect(location string) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Redirect(http.StatusTemporaryRedirect, location)
	}
}

func (ah *AppHandler) requireAuth(ctx *gin.Context) {
	ssn := sessions.Default(ctx)
	sid := ssn.Get("gnol-session-id")
	if sid != nil {
		gs := &storage.GnolSession{}

		if err := ah.dao.DB.Get(gs, "select * from gnol_session gs where session_id = $1 and gs.valid_until > now()", sid); err == nil {
			if gs != nil {
				gnolContext := NewGnolContext(gs)
				user, err := ah.dao.GetUser(gs.UserId)
				gnolContext.UserInfo = user
				if err != nil {
					panic(err)
				}
				ctx.Set("gnol-context", gnolContext)
				return
			}
		} else {
			log.Printf("Error while checking UserSession in DB: %v", err)
		}
	}

	ctx.Redirect(http.StatusTemporaryRedirect, "/users/login")
	ctx.Abort()
}

func getGnolContext(ctx *gin.Context) *GnolContext {
	return ctx.MustGet("gnol-context").(*GnolContext)
}

func userSessionMiddleware(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if session.Get("gnol-session-id") == nil {
		sid := xid.New().String()
		log.Printf("New Session %s", sid)
		session.Set("gnol-session-id", sid)
	}
	ctx.Next()
	if e := session.Save(); e != nil {
		panic(e)
	}
}

func (ah *AppHandler) initTemplates() {
	var allFiles []string
	var err error
	ah.templates = template.New("root")
	ah.templates = ah.templates.Funcs(template.FuncMap{"mod": mod, "inc": inc})
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

func (ah *AppHandler) renderTemplate(templateName string, ctx *gin.Context, renderCtx interface{}) {
	tpl, tlerr := ah.getTemplate(templateName)
	if tlerr != nil {
		renderError(tlerr, ctx.Writer)
	}
	re := tpl.Execute(ctx.Writer, renderCtx)
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

func inc(i int) int {
	return i + 1
}
