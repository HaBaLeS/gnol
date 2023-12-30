package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/HaBaLeS/gnol/server/cache"
	"github.com/HaBaLeS/gnol/server/jobs"
	"github.com/HaBaLeS/gnol/server/router"
	"github.com/HaBaLeS/gnol/server/storage"
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/HaBaLeS/go-logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	cfgPath := flag.String("c", "default.cfg", "Config File to use")
	flag.Parse()

	go startSwagger()

	gnol := NewServer(*cfgPath)
	gnol.Start()
}

func startSwagger() {
	r := gin.Default()

	r.Run(":8080")

}

// Application is the central struct connecting all submodules into one Application
// this struct supports the access between the modules.
type Application struct {
	Config     *util.ToolConfig
	HTTPServer *http.Server
	Handler    *router.AppHandler
	dao        *storage.DAO
	Logger     *logger.Logger
	Cache      *cache.ImageCache
	BGJobs     *jobs.JobRunner
}

// NewServer creates a new gnol Application
func NewServer(cfgPath string) *Application {
	log, err := logger.NewLogger()
	if err != nil {
		panic("Could not create Logger!")
	}
	cfg, err := util.ReadConfig(cfgPath)
	if err != nil {
		log.WarningF("%s not found using defaults", cfgPath)
	}

	a := &Application{
		Config: cfg,
		Logger: log,
	}

	log.InfoF("Using: http://%s:%d/users/login", a.Config.Hostname, a.Config.ListenPort)
	return a
}

// Start gnol, serve HTTP
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						gnol-token
// @description				your auth key for the API
func (a *Application) Start() {

	a.dao = storage.NewDAO(a.Config)

	a.Cache = cache.NewImageCache(a.Config)
	go a.Cache.RecoverCacheDir()

	a.BGJobs = jobs.NewJobRunner(a.dao, a.Config)
	a.BGJobs.StartMonitor()

	//TODO.md move router in server
	a.Handler = router.NewHandler(a.Config, a.Cache, a.BGJobs, a.dao)
	a.Handler.Routes()

	a.HTTPServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", a.Config.Hostname, a.Config.ListenPort),
		Handler: a.Handler.Router,
	}
	err := a.HTTPServer.ListenAndServe()
	if err != http.ErrServerClosed {
		panic(err)
	}
}

// Shutdown try's to end all modules gracefully where needed
func (a *Application) Shutdown() {
	a.BGJobs.StopMonitor()
	if a.HTTPServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := a.HTTPServer.Shutdown(ctx)
		if err != nil {
			panic(err)
		} else {
			a.HTTPServer = nil
		}
	}
}
