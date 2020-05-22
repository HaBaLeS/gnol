package server

import (
	"context"
	"fmt"
	"net/http"
	"playground.dahoam/util"
	"time"
)

type Session struct {
	config     *util.ToolConfig
	httpServer *http.Server
	handler    *AppHandler
	dao        *DAOHandler
}

func NewServer(cfgPath string) *Session {
	cfg, err := util.ReadConfig(cfgPath)
	if err != nil {
		panic(err) //FIXME don't panic exit graceful
	}

	s := &Session{
		config: cfg,
	}

	fmt.Printf("Using: http://%s:%d\n", s.config.ListenAddress, s.config.ListenPort)

	return s
}

func (s *Session) Start() {

	s.dao = NewDAO(*s)
	s.dao.Warmup()

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
