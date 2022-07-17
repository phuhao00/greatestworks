package main

import (
	"fmt"
	"greatestworks/network"
	"greatestworks/network/protocol/gen/player"
	"strconv"

	"google.golang.org/protobuf/proto"
)

type MessageHandler func(packet *network.ClientPacket)

type InputHandler func(param *InputParam)

//CreatePlayer 创建角色
func (c *Client) CreatePlayer(param *InputParam) {
	id := c.GetMessageIdByCmd(param.Command)

	if len(param.Param) != 2 {
		return
	}

	msg := &player.CSCreateUser{
		UserName: param.Param[0],
		Password: param.Param[1],
	}

	c.Transport(id, msg)
}

func (c *Client) OnCreatePlayerRsp(packet *network.ClientPacket) {
	fmt.Println("恭喜你创建角色成功")
}

func (c *Client) Login(param *InputParam) {
	id := c.GetMessageIdByCmd(param.Command)

	if len(param.Param) != 2 {
		return
	}

	msg := &player.CSLogin{
		UserName: param.Param[0],
		Password: param.Param[1],
	}

	c.Transport(id, msg)

}

func (c *Client) OnLoginRsp(packet *network.ClientPacket) {
	rsp := &player.SCLogin{}

	err := proto.Unmarshal(packet.Msg.Data, rsp)
	if err != nil {
		return
	}

	fmt.Println("登陆成功")
}

func (c *Client) AddFriend(param *InputParam) {
	id := c.GetMessageIdByCmd(param.Command)

	if len(param.Param) != 1 || len(param.Param[0]) == 0 { //""
		return
	}

	parseUint, err := strconv.ParseUint(param.Param[0], 10, 64)
	if err != nil {
		return
	}

	msg := &player.CSAddFriend{
		UId: parseUint,
	}
	c.Transport(id, msg)
}

func (c *Client) OnAddFriendRsp(packet *network.ClientPacket) {
	fmt.Println("add friend success !!")
}

func (c *Client) DelFriend(param *InputParam) {
	id := c.GetMessageIdByCmd(param.Command)

	if len(param.Param) != 1 || len(param.Param[0]) == 0 { //""
		return
	}

	parseUint, err := strconv.ParseUint(param.Param[0], 10, 64)
	if err != nil {
		return
	}

	msg := &player.CSDelFriend{
		UId: parseUint,
	}

	c.Transport(id, msg)
}

func (c *Client) OnDelFriendRsp(packet *network.ClientPacket) {
	fmt.Println("you have del friend success")

}

func (c *Client) SendChatMsg(param *InputParam) {
	id := c.GetMessageIdByCmd(param.Command)

	if len(param.Param) != 3 { //""
		return
	}

	parseUint, err := strconv.ParseUint(param.Param[0], 10, 64)
	if err != nil {
		return
	}
	parseInt32, err := strconv.ParseInt(param.Param[2], 10, 32)
	if err != nil {
		return
	}

	msg := &player.CSSendChatMsg{
		UId: parseUint,
		Msg: &player.ChatMessage{
			Content: param.Param[1],
			Extra:   nil,
		},
		Category: int32(parseInt32),
	}

	c.Transport(id, msg)
}

func (c *Client) OnSendChatMsgRsp(packet *network.ClientPacket) {
	fmt.Println("send  chat message success")

}
