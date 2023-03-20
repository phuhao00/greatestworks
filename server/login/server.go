package main

import (
	"fmt"
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
			BaseServer:  nil,
		}
		serverLogin.Initialize()
		var err error
		serverLogin.BaseServer, err = server.NewBaseServer(serverLogin.Conf.Me.Name, "")
		if err != nil {
			panic(fmt.Sprintf("[GetServer-initOnce] err:%v", err))
		}
	})

	return serverLogin
}

func (s *Server) Initialize() {
	//var consulInstance interface{}
	//consul 获取配置json串
	//todo load config
	s.Conf = config.Deserialize("")
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
