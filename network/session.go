package network

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type Session struct {
	UId     int64
	Conn    net.Conn
	IsClose bool
	packer  IPacker
	writeCh chan *Message
	//
}

func NewSession(conn net.Conn) *Session {
	return &Session{Conn: conn, packer: &NormalPacker{ByteOrder: binary.BigEndian}, writeCh: make(chan *Message, 1)}
}

func (s *Session) Run() {
	go s.Read()
	go s.Write()

}

func (s *Session) Read() {
	for {
		err := s.Conn.SetReadDeadline(time.Now().Add(time.Second))
		if err != nil {
			fmt.Println(err)
			continue
		}
		message, err := s.packer.Unpack(s.Conn)
		if _, ok := err.(net.Error); ok {
			continue
		}
		fmt.Println("receive message:", string(message.Data))
		s.writeCh <- &Message{
			ID:   555,
			Data: []byte("hi"),
		}
	}
}

func (s *Session) Write() {
	for {
		select {
		case resp := <-s.writeCh:
			s.send(resp)
		}
	}
}

func (s *Session) send(message *Message) {
	err := s.Conn.SetWriteDeadline(time.Now().Add(time.Second))
	if err != nil {
		fmt.Println(err)
		return
	}
	bytes, err := s.packer.Pack(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.Conn.Write(bytes)

}
