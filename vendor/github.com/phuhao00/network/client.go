package network

import (
	"github.com/phuhao00/spoor"
	"net"
	"runtime/debug"
	"sync/atomic"
)

type Client struct {
	*TcpConnX
	Address         string
	ChMsg           chan *Message
	OnMessageCb     func(message *Packet)
	logger          *spoor.Spoor
	bufferSize      int
	running         atomic.Value
	OnCloseCallBack func()
	closed          int32
}

func NewClient(address string, connBuffSize int, logger *spoor.Spoor) *Client {
	client := &Client{
		bufferSize: connBuffSize,
		Address:    address,
		logger:     logger,
		TcpConnX:   nil,
	}
	client.running.Store(false)
	return client
}

func (c *Client) Dial() (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", c.Address)

	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp6", nil, tcpAddr)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *Client) Run() {
	conn, err := c.Dial()
	if err != nil {
		c.logger.ErrorF("%v", err)
		return
	}
	tcpConnX, err := NewTcpConnX(conn, c.bufferSize, c.logger)
	if err != nil {
		c.logger.ErrorF("%v", err)
		return
	}
	c.TcpConnX = tcpConnX
	c.Impl = c
	c.Reset()
	c.running.Store(true)
	go c.Connect()
}

func (c *Client) OnClose() {
	if c.OnCloseCallBack != nil {
		c.OnCloseCallBack()
	}
	c.running.Store(false)
	c.TcpConnX.OnClose()
}

func (c *Client) OnMessage(data *Message, conn *TcpConnX) {

	c.Verify()

	defer func() {
		if err := recover(); err != nil {
			c.logger.ErrorF("[OnMessage] panic ", err, "\n", string(debug.Stack()))
		}
	}()

	c.OnMessageCb(&Packet{
		Msg:  data,
		Conn: conn,
	})
}

// Close 关闭连接
func (c *Client) Close() {
	if atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		c.Conn.Close()
		close(c.stopped)
	}
}
