package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"net"
	"time"
)

type NormalPacker struct {
	ByteOrder binary.ByteOrder
}

func (p *NormalPacker) Pack(msgID uint16, msg interface{}) ([]byte, error) {
	pbMsg, ok := msg.(proto.Message)
	if !ok {
		return []byte{}, fmt.Errorf("msg is not protobuf message")
	}
	data, err := proto.Marshal(pbMsg)
	if err != nil {
		return data, err
	}
	buffer := make([]byte, 8+8+len(data))
	p.ByteOrder.PutUint64(buffer[0:8], uint64(len(buffer)))
	p.ByteOrder.PutUint64(buffer[8:16], uint64(msgID))
	copy(buffer[16:], data)
	return buffer, nil
}

func (p *NormalPacker) Read(conn *TcpConnX) ([]byte, error) {
	err := conn.Conn.(*net.TCPConn).SetReadDeadline(time.Now().Add(time.Second * 10))
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, 8+8)

	if _, err := io.ReadFull(conn.Conn, buffer); err != nil {
		return nil, err
	}
	totalSize := p.ByteOrder.Uint64(buffer[:8])
	dataSize := totalSize - 8 - 8
	data := make([]byte, 8+dataSize)
	copy(data[:8], buffer[8:])
	if _, err := io.ReadFull(conn.Conn, data[8:]); err != nil {
		return nil, err
	}
	return data, nil
}

// Unpack id + protobuf data
func (p *NormalPacker) Unpack(data []byte) (*Message, error) {
	if len(data) < 2 {
		return nil, errors.New("protobuf data too short")
	}
	msgID := p.ByteOrder.Uint16(data[:2])
	msg := &Message{
		ID:   uint64(msgID),
		Data: data[2:],
	}
	return msg, nil
}
