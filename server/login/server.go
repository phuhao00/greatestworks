package main

import (
	"greatestworks/server"
	"greatestworks/server/login/config"
	"sync"
)

type Server struct {
	ProcessId   int
	httpHandler interface{}
	Timer       interface{}
	OPenTime    int64
	Conf        *config.Config
	*server.BaseServer
}

var (
	serverLogin *Server
	initOnce    sync.Once
)

func GetServer() *Server {
	initOnce.Do(func() {
		serverLogin = &Server{
			ProcessId:   0,
			httpHandler: nil,
			Timer:       nil,
			OPenTime:    0,
			Conf:        nil,
		}
	})
	return serverLogin
}

func (s *Server) Initialize(conf *config.Config) {

}

func (s *Server) RegisterTimer() {

}

func (s *Server) Run() {
	for {
		select {}
	}
}

func (s *Server) ServiceRegister() {
	//consul register
}

func (s *Server) UpdateRegister() {
	//update consul register
}

func (s *Server) GetOtherService() {

}
