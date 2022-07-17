package network

import "net"

type ClientPacket struct {
	Msg  *Message
	Conn net.Conn
}

type SessionPacket struct {
	Msg  *Message
	Sess *Session
}
