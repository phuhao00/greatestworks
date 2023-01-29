package network

import (
	"encoding/binary"
)

const (
	cheapPrependSize = 8
	initialSize      = 1024
)

// ByteBuffer 字节buff
type ByteBuffer struct {
	mBuffer             []byte
	mCapacity           int32
	readIndex           int32
	writeIndex          int32
	reservedPrependSize int32
	littleEndian        bool
}

// NewByteBuffer 创建一个字节buffer
func NewByteBuffer() *ByteBuffer {
	return &ByteBuffer{
		mBuffer:             make([]byte, cheapPrependSize+initialSize),
		mCapacity:           cheapPrependSize + initialSize,
		readIndex:           cheapPrependSize,
		writeIndex:          cheapPrependSize,
		reservedPrependSize: cheapPrependSize,
		littleEndian:        true,
	}
}

// SetByteOrder It's dangerous to call the method on reading or writing
func (bf *ByteBuffer) SetByteOrder(littleEndian bool) {
	bf.littleEndian = littleEndian
}

// Length ...
func (bf *ByteBuffer) Length() int32 {
	return bf.writeIndex - bf.readIndex
}

// Swap ...
func (bf *ByteBuffer) Swap(other *ByteBuffer) {
}

// Skip advances the reading index of the buffer
func (bf *ByteBuffer) Skip(len int32) {
	if len < bf.Length() {
		bf.readIndex = bf.readIndex + len
	} else {
		bf.Reset()
	}
}

// Retrieve ...
func (bf *ByteBuffer) Retrieve(len int32) {
	bf.Skip(len)
}

// Reset ...
func (bf *ByteBuffer) Reset() {
	bf.Truncate(0)
}

// Truncate ...
func (bf *ByteBuffer) Truncate(n int32) {
	if n == 0 {
		bf.readIndex = bf.reservedPrependSize
		bf.writeIndex = bf.reservedPrependSize
	} else if bf.writeIndex > (bf.readIndex + n) {
		bf.writeIndex = bf.readIndex + n
	}
}

// Reserve ...
func (bf *ByteBuffer) Reserve(len int32) {
	if bf.mCapacity >= len+bf.reservedPrependSize {
		return
	}
	bf.grow(len + bf.reservedPrependSize)
}

// Append ...
func (bf *ByteBuffer) Append(buff []byte) {
	size := len(buff)
	if size == 0 {
		return
	}
	bf.write(buff, int32(size))
}

// AppendInt64 ...
func (bf *ByteBuffer) AppendInt64(x int64) {
	buff := make([]byte, 8)
	if bf.littleEndian {
		binary.LittleEndian.PutUint64(buff, uint64(x))
	} else {
		binary.BigEndian.PutUint64(buff, uint64(x))
	}
	bf.write(buff, 8)
}

// AppendInt32 ...
func (bf *ByteBuffer) AppendInt32(x int32) {
	buff := make([]byte, 4)
	if bf.littleEndian {
		binary.LittleEndian.PutUint32(buff, uint32(x))
	} else {
		binary.BigEndian.PutUint32(buff, uint32(x))
	}
	bf.write(buff, 4)
}

// AppendInt16 ...
func (bf *ByteBuffer) AppendInt16(x int16) {
	buff := make([]byte, 2)
	if bf.littleEndian {
		binary.LittleEndian.PutUint16(buff, uint16(x))
	} else {
		binary.BigEndian.PutUint16(buff, uint16(x))
	}
	bf.write(buff, 2)
}

// ReadInt64 ...
func (bf *ByteBuffer) ReadInt64() int64 {
	buff := bf.mBuffer[bf.readIndex : bf.readIndex+8]
	var result uint64
	if bf.littleEndian {
		result = binary.LittleEndian.Uint64(buff)
	} else {
		result = binary.BigEndian.Uint64(buff)
	}
	bf.Skip(8)
	return int64(result)
}

// ReadInt32 ...
func (bf *ByteBuffer) ReadInt32() int32 {
	buff := bf.mBuffer[bf.readIndex : bf.readIndex+4]
	var result uint32
	if bf.littleEndian {
		result = binary.LittleEndian.Uint32(buff)
	} else {
		result = binary.BigEndian.Uint32(buff)
	}
	bf.Skip(4)
	return int32(result)
}

// ReadInt16 ...
func (bf *ByteBuffer) ReadInt16() int16 {
	buff := bf.mBuffer[bf.readIndex : bf.readIndex+2]
	var result uint16
	if bf.littleEndian {
		result = binary.LittleEndian.Uint16(buff)
	} else {
		result = binary.BigEndian.Uint16(buff)
	}
	bf.Skip(2)
	return int16(result)
}

// NextBytes 读取N字节
func (bf *ByteBuffer) NextBytes(len int32) []byte {
	msgData := bf.mBuffer[bf.readIndex : bf.readIndex+len]
	bf.readIndex += len
	return msgData
}

// EnsureWritableBytes ...
func (bf *ByteBuffer) EnsureWritableBytes(len int32) {
	if bf.writableBytes() < len {
		bf.grow(len)
	}
}

func (bf *ByteBuffer) grow(len int32) {
	if bf.writableBytes()+bf.prependableBytes() < len+bf.reservedPrependSize {
		newCap := (bf.mCapacity << 1) + len
		buff := make([]byte, newCap)
		copy(buff, bf.mBuffer)
		bf.mCapacity = newCap
		bf.mBuffer = buff
	} else {
		readable := bf.Length()
		copy(bf.mBuffer[bf.reservedPrependSize:], bf.mBuffer[bf.readIndex:bf.writeIndex])
		bf.readIndex = bf.reservedPrependSize
		bf.writeIndex = bf.readIndex + readable
	}
}

func (bf *ByteBuffer) write(buff []byte, len int32) {
	bf.EnsureWritableBytes(len)
	copy(bf.mBuffer[bf.writeIndex:], buff)
	bf.writeIndex = bf.writeIndex + len
}

func (bf *ByteBuffer) writableBytes() int32 {
	return bf.mCapacity - bf.writeIndex
}

func (bf *ByteBuffer) prependableBytes() int32 {
	return bf.readIndex
}

// WriteBuff ...
func (bf *ByteBuffer) WriteBuff() []byte {
	buffLen := int32(len(bf.mBuffer))
	if bf.writeIndex >= buffLen {
		return nil
	}
	return bf.mBuffer[bf.writeIndex:]
}

// ReadBuff ...
func (bf *ByteBuffer) ReadBuff() []byte {
	buffLen := int32(len(bf.mBuffer))
	if bf.readIndex >= buffLen {
		return nil
	}
	return bf.mBuffer[bf.readIndex:]
}

// WriteBytes 写入n字节
func (bf *ByteBuffer) WriteBytes(n int32) {
	bf.writeIndex += n
}
