package network

import (
	"fmt"
	"net"
)

type Server struct {
	tcpListener     net.Listener
	OnSessionPacket func(packet *SessionPacket)
	Address         string
}

func NewServer(address string) *Server {

	s := &Server{Address: address}

	return s

}

func (s *Server) Run() {
	resolveTCPAddr, err := net.ResolveTCPAddr("tcp6", s.Address)
	if err != nil {
		panic(err)
	}
	tcpListener, err := net.ListenTCP("tcp6", resolveTCPAddr)
	if err != nil {
		panic(err)
	}
	s.tcpListener = tcpListener
	for {
		conn, err := s.tcpListener.Accept()
		if err != nil {
			if _, ok := err.(net.Error); ok {
				fmt.Println(err)
				continue
			}
		}

		newSession := NewSession(conn)
		newSession.MessageHandler = s.OnSessionPacket
		SessionMgrInstance.AddSession(newSession)
		newSession.Run()
		SessionMgrInstance.DelSession(newSession.UId)
	}
}
