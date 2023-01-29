package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"math"
)

type BufferPacker struct {
	lenMsgLen int32
	minMsgLen uint32
	maxMsgLen uint32
	recvBuff  *ByteBuffer
	sendBuff  *ByteBuffer
	byteOrder binary.ByteOrder
}

func newInActionPacker() *BufferPacker {
	msgParser := &BufferPacker{
		lenMsgLen: 4,
		minMsgLen: 2,
		maxMsgLen: 2 * 1024 * 1024,
		recvBuff:  NewByteBuffer(),
		sendBuff:  NewByteBuffer(),
		byteOrder: binary.LittleEndian,
	}
	return msgParser
}

// SetMsgLen It's dangerous to call the method on reading or writing
func (p *BufferPacker) SetMsgLen(lenMsgLen int32, minMsgLen uint32, maxMsgLen uint32) {
	if lenMsgLen == 1 || lenMsgLen == 2 || lenMsgLen == 4 {
		p.lenMsgLen = lenMsgLen
	}
	if minMsgLen != 0 {
		p.minMsgLen = minMsgLen
	}
	if maxMsgLen != 0 {
		p.maxMsgLen = maxMsgLen
	}

	var max uint32
	switch p.lenMsgLen {
	case 1:
		max = math.MaxUint8
	case 2:
		max = math.MaxUint16
	case 4:
		max = math.MaxUint32
	}
	if p.minMsgLen > max {
		p.minMsgLen = max
	}
	if p.maxMsgLen > max {
		p.maxMsgLen = max
	}
}

// Read goroutine safe
func (p *BufferPacker) Read(conn *TcpConnX) ([]byte, error) {

	p.recvBuff.EnsureWritableBytes(p.lenMsgLen)

	readLen, err := io.ReadFull(conn, p.recvBuff.WriteBuff()[:p.lenMsgLen])
	// read len
	if err != nil {
		return nil, fmt.Errorf("%v readLen:%v", err, readLen)
	}
	p.recvBuff.WriteBytes(int32(readLen))

	// parse len
	var msgLen uint32
	switch p.lenMsgLen {
	case 2:
		msgLen = uint32(p.recvBuff.ReadInt16())
	case 4:
		msgLen = uint32(p.recvBuff.ReadInt32())
	}

	// check len
	if msgLen > p.maxMsgLen {
		return nil, errors.New("message too long")
	} else if msgLen < p.minMsgLen {
		return nil, errors.New("message too short")
	}

	p.recvBuff.EnsureWritableBytes(int32(msgLen))

	rLen, err := io.ReadFull(conn, p.recvBuff.WriteBuff()[:msgLen])
	if err != nil {
		return nil, fmt.Errorf("%v msgLen:%v readLen:%v", err, msgLen, rLen)
	}
	p.recvBuff.WriteBytes(int32(rLen))

	/*
		// 保留了2字节flag 暂时未处理
		var flag uint16
		flag = uint16(p.recvBuff.ReadInt16())
	*/
	p.recvBuff.Skip(2) // 跳过2字节保留字段

	// 减去2字节的保留字段长度
	return p.recvBuff.NextBytes(int32(msgLen - 2)), nil

}

// goroutine safe
func (p *BufferPacker) Write(conn *TcpConnX, buff ...byte) error {
	// get len
	msgLen := uint32(len(buff))

	// check len
	if msgLen > p.maxMsgLen {
		return errors.New("message too long")
	} else if msgLen < p.minMsgLen {
		return errors.New("message too short")
	}

	// write len
	switch p.lenMsgLen {
	case 2:
		p.sendBuff.AppendInt16(int16(msgLen))
	case 4:
		p.sendBuff.AppendInt32(int32(msgLen))
	}

	p.sendBuff.Append(buff)
	// write data
	writeBuff := p.sendBuff.ReadBuff()[:p.sendBuff.Length()]

	_, err := conn.Write(writeBuff)

	p.sendBuff.Reset()

	return err
}

func (p *BufferPacker) reset() {
	p.recvBuff = NewByteBuffer()
	p.sendBuff = NewByteBuffer()
}

func (p *BufferPacker) Pack(msgID uint16, msg interface{}) ([]byte, error) {
	pbMsg, ok := msg.(proto.Message)
	if !ok {
		return []byte{}, fmt.Errorf("msg is not protobuf message")
	}
	// data
	data, err := proto.Marshal(pbMsg)
	if err != nil {
		return data, err
	}
	// 4byte = len(flag)[2byte] + len(msgID)[2byte]
	buf := make([]byte, 4+len(data))
	if p.byteOrder == binary.LittleEndian {
		binary.LittleEndian.PutUint16(buf[0:2], 0)
		binary.LittleEndian.PutUint16(buf[2:], msgID)
	} else {
		binary.BigEndian.PutUint16(buf[0:2], 0)
		binary.BigEndian.PutUint16(buf[2:], msgID)
	}
	copy(buf[4:], data)
	return buf, err
}

// Unpack id + protobuf data
func (p *BufferPacker) Unpack(data []byte) (*Message, error) {
	if len(data) < 2 {
		return nil, errors.New("protobuf data too short")
	}
	msgID := p.byteOrder.Uint16(data[:2])
	msg := &Message{
		ID:   uint64(msgID),
		Data: data[2:],
	}
	return msg, nil
}
