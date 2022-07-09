package network

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type Client struct {
	Address string
	packer  IPacker
	chMsg   chan *Message
}

func NewClient(address string) *Client {
	return &Client{
		Address: address,
		packer: &NormalPacker{
			ByteOrder: binary.BigEndian,
		},
		chMsg: make(chan *Message, 1),
	}
}

func (c *Client) Run() {
	conn, err := net.Dial("tcp6", c.Address)
	if err != nil {
		fmt.Println(err)
		return
	}
	go c.Read(conn)
	go c.Write(conn)

}

func (c *Client) Write(conn net.Conn) {
	tick := time.NewTicker(time.Second)
	for {
		select {
		case <-tick.C:
			c.chMsg <- &Message{
				ID:   111,
				Data: []byte("hello world "),
			}
		case msg := <-c.chMsg:
			c.Send(conn, msg)
		}
	}
}

func (c *Client) Send(conn net.Conn, message *Message) {
	pack, err := c.packer.Pack(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.Write(pack)
}

func (c *Client) Read(conn net.Conn) {
	for {
		message, err := c.packer.Unpack(conn)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("read msg:", string(message.Data))
	}
}
