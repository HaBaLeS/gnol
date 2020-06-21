package server

import (
	"context"
	"fmt"
	"github.com/HaBaLeS/gnol/util"
	"github.com/HaBaLeS/go-logger"
	"net/http"
	"time"
)

type Session struct {
	config     *util.ToolConfig
	httpServer *http.Server
	handler    *AppHandler
	dao        *DAOHandler
	logger     *logger.Logger
	cache      *ImageCache
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
		config: cfg,
		logger: log,
	}

	log.InfoF("Using: http://%s:%d/comics", s.config.ListenAddress, s.config.ListenPort)
	return s
}

func (s *Session) Start() {

	s.dao = NewDAO(*s)
	s.dao.Warmup()

	s.cache = NewImageCache(s)

	//TODO move router in server
	s.handler = NewHandler(s)
	s.handler.SetupRoutes()

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.config.ListenAddress, s.config.ListenPort),
		Handler: s.handler.router,
	}
	err := s.httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		panic(err)
	}
}

func (s *Session) Shutdown() {
	if s.httpServer != nil {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			panic(err)
		} else {
			s.httpServer = nil
		}
	}
}
