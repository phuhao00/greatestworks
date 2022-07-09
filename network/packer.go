package network

import (
	"encoding/binary"
	"io"
	"net"
	"time"
)

type NormalPacker struct {
	ByteOrder binary.ByteOrder
}

func (p *NormalPacker) Pack(message *Message) ([]byte, error) {
	buffer := make([]byte, 8+8+len(message.Data))
	p.ByteOrder.PutUint64(buffer[0:8], uint64(len(buffer)))
	p.ByteOrder.PutUint64(buffer[8:16], message.ID)
	copy(buffer[16:], message.Data)
	return buffer, nil
}

func (p *NormalPacker) Unpack(reader io.Reader) (*Message, error) {
	err := reader.(*net.TCPConn).SetReadDeadline(time.Now().Add(time.Second * 10))
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, 8+8)

	if _, err := io.ReadFull(reader, buffer); err != nil {
		return nil, err
	}
	totalSize := p.ByteOrder.Uint64(buffer[:8])
	Id := p.ByteOrder.Uint64(buffer[8:])
	dataSize := totalSize - 8 - 8
	data := make([]byte, dataSize)

	if _, err := io.ReadFull(reader, data); err != nil {
		return nil, err
	}
	msg := &Message{
		ID:   Id,
		Data: data,
	}
	return msg, nil
}
