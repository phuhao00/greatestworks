package client

import (
	"github.com/phuhao00/fuse"
	"github.com/phuhao00/greatestworks-proto/messageId"
	"github.com/phuhao00/greatestworks-proto/server_common"
	"github.com/phuhao00/network"
	"greatestworks/server/world/server"
	"sync"
)

type Client struct {
	tcpClient      *network.TcpClient
	gameMsgReg     sync.Once
	serverID       string
	serverAddr     string
	handleMessages []uint16
	Weights        int32
	Index          int
	LogicRouter    *fuse.LogicRouter
}

func (c *Client) connect() {
	c.tcpClient.Connect()
	c.gameMsgReg.Do(c.registerMsgHandler)
	c.registerToGameServer()
}

func (c *Client) checkConnect() {
	//if c.tcpClient.IsRunning() {
	//	return
	//}
	c.connect()
}

func (c *Client) registerMsgHandler() {
	c.LogicRouter.Register(uint64(messageId.MessageId_RegisterToGameRet), c.registerToGameHandler)
	c.LogicRouter.Register(uint64(messageId.MessageId_CreateGameRet), c.createGameHandler)
	c.LogicRouter.Register(uint64(messageId.MessageId_SrvMsgToOnlinePlayer), c.gameMsgToOnlinePlayerHandler)
}

func (c *Client) sendMsg(messageId messageId.MessageId, msg interface{}) bool {
	c.tcpClient.AsyncSend(uint64(messageId), msg)
	return true
}

func (c *Client) registerToGameServer() {
	regMsg := &server_common.RegisterToGame{
		ServerID:   server.Oasis.Id,
		ServerType: uint32(server_common.ServerCategory_World),
		ServerAddr: server.Oasis.Config.HttpAddress,
		Token:      "greatest_works",
		ScrSerType: uint32(server_common.ServerCategory_World),
		SerRpcAddr: server.Oasis.Config.RpcAddr,
	}
	c.sendMsg(messageId.MessageId_RegisterToGame, regMsg)
}
