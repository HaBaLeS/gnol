package server

import (
	"context"
	"net/http"
	"playground.dahoam/util"
	"time"
)

type Session struct {
	config *util.ToolConfig
	httpServer *http.Server
	handler *AppHandler
}

func NewServer() *Session{
	s := &Session{
		//... init something ?
	}
	return s
}

func (s *Session) Start(){

	cfg, err := util.ReadConfig("default.cfg")
	if err != nil {
		panic(err)
	}
	s.config = cfg

	//TODO move router in server
	s.handler = NewHandler(s)
	s.handler.SetupRoutes()

	s.httpServer = &http.Server{
		Addr: "192.168.1.248:6969",
		Handler: s.handler.router,
	}
	err = s.httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		panic(err)
	}
}

func (s *Session) Shutdown(){
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


