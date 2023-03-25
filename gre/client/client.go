package main

import (
	"github.com/phuhao00/greatestworks-proto/messageId"
	"greatestworks/aop/logger"
	"os"
	"syscall"

	"github.com/phuhao00/network"
)

type Client struct {
	cli             *network.Client
	inputHandlers   map[string]InputHandler
	messageHandlers map[messageId.MessageId]MessageHandler
	console         *ClientConsole
	chInput         chan *InputParam
}

func NewClient() *Client {
	c := &Client{
		cli:             network.NewClient(":8023", 200, logger.Logger),
		inputHandlers:   map[string]InputHandler{},
		messageHandlers: map[messageId.MessageId]MessageHandler{},
		console:         NewClientConsole(),
	}
	c.cli.OnMessageCb = c.OnMessage
	c.cli.ChMsg = make(chan *network.Message, 1)
	c.chInput = make(chan *InputParam, 1)
	c.console.chInput = c.chInput
	//https://github.com/phuhao00/greatestworks-proto.git
	//github.com/phuhao00/greatestworks-proto
	return c
}

func (c *Client) Run() {
	go func() {
		for {
			select {
			case input := <-c.chInput:
				inputHandler := c.inputHandlers[input.Command]
				if inputHandler != nil {
					inputHandler(input)
				}
			}
		}
	}()
	go c.console.Run()
	go c.cli.Run()
}

func (c *Client) OnMessage(packet *network.Packet) {
	if handler, ok := c.messageHandlers[messageId.MessageId(packet.Msg.ID)]; ok {
		handler(packet)
	}
}

func (c *Client) OnSystemSignal(signal os.Signal) bool {
	logger.Logger.InfoF("[Client] 收到信号 %v \n", signal)
	tag := true
	switch signal {
	case syscall.SIGHUP:
		//todo
	case syscall.SIGPIPE:
	default:
		logger.Logger.InfoF("[Client] 收到信号准备退出...")
		tag = false

	}
	return tag
}
