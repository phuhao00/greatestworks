package gateway

import (
	"github.com/phuhao00/fuse"
	"github.com/phuhao00/greatestworks-proto/gateway"
	"github.com/phuhao00/greatestworks-proto/messageId"
	"github.com/phuhao00/greatestworks-proto/server_common"
	"github.com/phuhao00/network"
	"greatestworks/server/world/config"
	"greatestworks/server/world/server"
	"sync"
)

type Client struct {
	tcpClient             *network.TcpClient
	msgId2HandlerRegister sync.Once
	serverID              string
	serverAddr            string
	msgID                 []uint32
	ZoneId                int
	LogicRouter           fuse.LogicRouter
}

func (c *Client) connect() {
	c.tcpClient.Connect()
	c.msgId2HandlerRegister.Do(c.registerMsgHandler)
	c.registerToGateway()
}

func (c *Client) checkConnect() {
	if !c.tcpClient.IsRunning() {
		return
	}
	c.connect()
}

func (c *Client) registerMsgHandler() {
	c.LogicRouter.Register(uint64(messageId.MessageId_RegisterToGatewayRet), c.registerToGatewayHandler)
	c.LogicRouter.Register(uint64(messageId.MessageId_UpdateOnlineInfoRet), c.updateOnlineInfoRetHandler)
	c.LogicRouter.Register(uint64(messageId.MessageId_ClientOnline), c.clientOnlineHandler)
	c.LogicRouter.Register(uint64(messageId.MessageId_ClientOffline), c.clientOfflineHandler)
	c.LogicRouter.Register(uint64(messageId.MessageId_GatewayForwardPacket), c.gatewayForwardPacketHandler)
}

func (c *Client) registerToGateway() {
	regMsg := &gateway.RegisterToGateway{
		ServerID:   server.Oasis.Id,
		ServerType: uint32(server_common.ServerCategory_World),
		ServerAddr: server.Oasis.Config.HttpAddress,
		Token:      "greatest_works",
		Msgids:     c.msgID,
		ProcIndex:  uint32(server.Oasis.Pid),
		MaxPlayer:  1000,
		ZoneId:     int32(c.ZoneId),
	}

	if conf := config.GetServerList(); conf != nil {
		if server.Oasis.Config.Global.ServerType == 0 {
			regMsg.Name = conf.Name
		} else if server.Oasis.Config.Global.ServerType == 1 {
			regMsg.Name = conf.Name1
		} else {
			regMsg.Name = conf.Name
		}

		var maxPlayerNum int32
		if server.Oasis.Config.MaxPlayerNum != 0 {
			maxPlayerNum = server.Oasis.Config.MaxPlayerNum
		} else {
			serverCfg := config.GetServerList()
			if serverCfg == nil {
				return
			}
			maxPlayerNum = serverCfg.MaxPlayer
		}

		regMsg.MaxPlayer = maxPlayerNum
	}
	c.sendMsg(messageId.MessageId_RegisterToGateway, regMsg)
}

func (c *Client) sendMsg(messageId messageId.MessageId, msg interface{}) bool {
	c.tcpClient.AsyncSend(uint64(messageId), msg)
	return true
}
