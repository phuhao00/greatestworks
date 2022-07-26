package main

import (
	"github.com/phuhao00/network"
	"greatestworks/network/protocol/gen/messageId"

	"google.golang.org/protobuf/proto"
)

func (c *Client) InputHandlerRegister() {
	c.inputHandlers[messageId.MessageId_CSCreatePlayer.String()] = c.CreatePlayer
	c.inputHandlers[messageId.MessageId_CSLogin.String()] = c.Login
	c.inputHandlers[messageId.MessageId_CSAddFriend.String()] = c.AddFriend
	c.inputHandlers[messageId.MessageId_CSDelFriend.String()] = c.DelFriend
	c.inputHandlers[messageId.MessageId_CSSendChatMsg.String()] = c.SendChatMsg
}

func (c *Client) GetMessageIdByCmd(cmd string) messageId.MessageId {
	mid, ok := messageId.MessageId_value[cmd]
	if ok {
		return messageId.MessageId(mid)
	}
	return messageId.MessageId_None
}

func (c *Client) Transport(id messageId.MessageId, message proto.Message) {
	bytes, err := proto.Marshal(message)
	if err != nil {
		return
	}
	c.cli.ChMsg <- &network.Message{
		ID:   uint64(id),
		Data: bytes,
	}
}
