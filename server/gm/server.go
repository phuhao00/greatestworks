package main

import (
	"github.com/golang/protobuf/proto"
	"greatestworks/server"
)

type Server struct {
	*server.BaseServer
	//todo gateway client
	toGateWay chan proto.Message
	//todo  mysql client 维护gm 账户
}

func (s *Server) OnGateWayMsg(message proto.Message) {
	//todo send to gateway
}
