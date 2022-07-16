package network

import "net"

type Server struct {
	tcpListener     net.Listener
	OnSessionPacket func(packet *SessionPacket)
}

func NewServer(address string) *Server {
	resolveTCPAddr, err := net.ResolveTCPAddr("tcp6", address)
	if err != nil {
		panic(err)
	}
	tcpListener, err := net.ListenTCP("tcp6", resolveTCPAddr)
	if err != nil {
		panic(err)
	}
	s := &Server{}
	s.tcpListener = tcpListener
	return s

}

func (s *Server) Run() {
	for {
		conn, err := s.tcpListener.Accept()
		if err != nil {
			if _, ok := err.(net.Error); ok {
				continue
			}
		}
		go func() {
			newSession := NewSession(conn)
			SessionMgrInstance.AddSession(newSession)
			newSession.Run()
			SessionMgrInstance.DelSession(newSession.UId)
		}()
	}
}
