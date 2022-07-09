package network

import "io"

type IPacker interface {
	Pack(message *Message) ([]byte, error)
	Unpack(reader io.Reader) (*Message, error)
}
