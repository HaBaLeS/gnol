package session

import (
	"context"
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

type Session struct {
	Config     *util.ToolConfig
	HttpServer *http.Server
	Handler    *router.AppHandler
	Dao        *dao.DAOHandler
	Logger     *logger.Logger
	Cache      *cache.ImageCache
	BGJobs 	   *conversion.JobRunner
}

func NewServer(cfgPath string) *Session {
	log, err := logger.NewLogger()
	if err != nil {
		panic("Could not create Logger!")
	}
	cfg, err := util.ReadConfig(cfgPath)
	if err != nil {
		log.WarningF("%s not found using defaults", cfgPath)
	}

	s := &Session{
		Config: cfg,
		Logger: log,
	}

	log.InfoF("Using: http://%s:%d/comics", s.Config.ListenAddress, s.Config.ListenPort)
	return s
}

func (s *Session) Start() {



	s.Dao = dao.NewDAO(s.Logger, s.Config)
	s.Dao.Warmup()

	s.Cache = cache.NewImageCache(s.Config)
	go s.Cache.RecoverCacheDir()

	s.BGJobs = conversion.NewJobRunner(s.Config.JobDirectory, s.Dao)
	s.BGJobs.StartMonitor()

	//TODO move router in server
	s.Handler = router.NewHandler(s.Config, s.Dao, s.Cache,s.BGJobs)
	s.Handler.SetupRoutes()
	s.Handler.SetupUploads()

	s.HttpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.Config.ListenAddress, s.Config.ListenPort),
		Handler: s.Handler.Router,
	}
	err := s.HttpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		panic(err)
	}
}

func (s *Session) Shutdown() {
	s.BGJobs.StopMonitor()
	if s.HttpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := s.HttpServer.Shutdown(ctx)
		if err != nil {
			panic(err)
		} else {
			s.HttpServer = nil
		}
	}
}
