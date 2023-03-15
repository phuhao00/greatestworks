package client

import "github.com/phuhao00/network"

type Server struct {
	real *network.Server
}

func (s *Server) Loop() {

	for {
		select {
		//impl Message

		}
	}
}
