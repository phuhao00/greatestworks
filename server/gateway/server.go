package main

import (
	"github.com/phuhao00/network"
	"greatestworks/server"
)

type Server struct {
	real *network.Server
	*server.BaseServer
}

func (s *Server) Loop() {

	for {
		select {
		//impl Message

		}
	}
}

func (s *Server) KickUser() {

}

func (s *Server) KickAllUser() {

}

func (s *Server) TransMessageToGateway() {

}

func (s *Server) UpdateRegister() {

}

func (s *Server) CheckRegister() bool {
	return true
}

func (s *Server) ReLoginAll() {

}
