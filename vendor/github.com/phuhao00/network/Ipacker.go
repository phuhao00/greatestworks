package network

type IPacker interface {
	Pack(msgID uint16, msg interface{}) ([]byte, error)
	Read(*TcpConnX) ([]byte, error)
	Unpack([]byte) (*Message, error)
}
