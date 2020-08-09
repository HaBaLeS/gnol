package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/HaBaLeS/gnol/server/cache"
	"github.com/HaBaLeS/gnol/server/conversion"
	"github.com/HaBaLeS/gnol/server/dao"
	"github.com/HaBaLeS/gnol/server/router"
	"github.com/HaBaLeS/gnol/server/util"
	"github.com/HaBaLeS/go-logger"
	"net/http"
	"time"
)

//go:generate go run -tags=dev gen.go

func main() {
	cfgPath := flag.String("c", "default.cfg", "Config File to use")
	flag.Parse()

	gnol := NewServer(*cfgPath)
	gnol.Start()
}

//Application is the central struct connecting all submodules into one Application
//this struct supports the access between the modules.
type Application struct {
	Config     *util.ToolConfig
	HTTPServer *http.Server
	Handler    *router.AppHandler
	Dao        *dao.DAOHandler
	Logger     *logger.Logger
	Cache      *cache.ImageCache
	BGJobs     *conversion.JobRunner
}

//NewServer creates a new gnol Application
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

	log.InfoF("Using: http://%s:%d/comics", a.Config.ListenAddress, a.Config.ListenPort)
	return a
}

//Start gnol, serve HTTP
func (a *Application) Start() {

	a.Dao = dao.NewDAO(a.Logger, a.Config)

	a.Cache = cache.NewImageCache(a.Config)
	go a.Cache.RecoverCacheDir()

	a.BGJobs = conversion.NewJobRunner(a.Dao)
	a.BGJobs.StartMonitor()

	//TODO move router in server
	a.Handler = router.NewHandler(a.Config, a.Dao, a.Cache, a.BGJobs)
	a.Handler.SetupRoutes()
	a.Handler.SetupUploads()
	a.Handler.SetupUserRouting()

	a.HTTPServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", a.Config.ListenAddress, a.Config.ListenPort),
		Handler: a.Handler.Router,
	}
	err := a.HTTPServer.ListenAndServe()
	if err != http.ErrServerClosed {
		panic(err)
	}
}

//Shutdown try's to end all modules gracefully where needed
func (a *Application) Shutdown() {
	a.Dao.Close()
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
