package gateway

import (
	"github.com/phuhao00/greatestworks-proto/gateway"
	"github.com/phuhao00/greatestworks-proto/messageId"
	"github.com/phuhao00/greatestworks-proto/player"
	"github.com/phuhao00/greatestworks-proto/server_common"
	"greatestworks/aop/logger"
	"greatestworks/server/world/server"
)

// registerToGatewayHandler ...
func (c *Client) registerToGatewayHandler(msgID uint64, data []byte) {
	msg := &gateway.RegisterToGatewayRet{}
	err := c.LogicRouter.Unmarshal(data, msg)
	if err != nil {
		logger.Error("[registerToGatewayHandler] receive data:%v msg:%v", data, msg)
		return
	}

	if msg.Result == uint32(server_common.ServerRegisterStatus_OK) {
		logger.Info("[registerToGatewayHandler] serverID:%v address:%v", c.serverID, c.serverAddr)
	}
}

// updateOnlineInfoRetHandler ...
func (c *Client) updateOnlineInfoRetHandler(msgID uint64, data []byte) {
	msg := &gateway.UpdateOnlineInfoRet{}
	err := c.LogicRouter.Unmarshal(data, msg)
	if err != nil {
		logger.Error("[updateOnlineInfoRetHandler] receive data:%v msg:%v", data, msg)
		return
	}
	if msg.Result == uint32(server_common.ServerRegisterStatus_OK) {
		logger.Info("[updateOnlineInfoRetHandler] serverID:%v Addr:%v", c.serverID, c.serverAddr)
	}
}

// gatewayForwardPacketHandler ...
func (c *Client) gatewayForwardPacketHandler(msgID uint64, data []byte) {
	msg := &gateway.GatewayForwardPacket{}
	err := c.LogicRouter.Unmarshal(data, msg)
	if err != nil {
		logger.Error("[gatewayForwardPacketHandler] receive data:%v msg:%v", data, msg)
		return
	}
	c.msgPacketHandle(msg.Userid, msg.Data)
}

// msgPacketHandle ...
func (c *Client) msgPacketHandle(userID uint64, data []byte) {
	msgID, err := c.LogicRouter.Route(data)
	if err != nil {
		err = server.Oasis.HandlePlayersMsgPacket(userID, data)
	}
	if err != nil {
		logger.Error("[msgPacketHandle] route message:%v error: %v", msgID, err)
	}
}

// clientOnlineHandler ...
func (c *Client) clientOnlineHandler(msgID uint64, data []byte) {
	msg := &gateway.ClientOnline{}
	err := c.LogicRouter.Unmarshal(data, msg)
	if err != nil {
		logger.Error("[clientOnlineHandler] receive data:%v msg:%v", data, msg)
		return
	}

	if msg.Reconnected {
		msgSend := &gateway.ClientOnlineRet{
			Result:       1,
			Userid:       msg.Userid,
			Reconnection: true,
			Players:      int32(server.Oasis.GetPlayersNum()),
		}
		p := server.Oasis.GetPlayer(msg.Userid)
		if p == nil {
			msgSend.Result = 0
		}
		c.sendMsg(messageId.MessageId_ClientOnlineRet, msgSend)
		return
	}
	if msg.Userid == 0 {
		return
	}
	onlineData := &player.PlayerData{
		PlayerID:  msg.Userid,
		GatewayID: c.serverID,
		RemoteIp:  msg.RemoteIp,
		GatewayIp: msg.GatewayIp,
		Version:   msg.Version,
	}
	server.Oasis.ChanPlayerOffline <- onlineData
	logger.Info("[clientOnlineHandler] Userid:%v ,chanPlayerOnline length: %v", msg.Userid, len(server.Oasis.ChanPlayerOffline))
}

// clientOfflineHandler ...
func (c *Client) clientOfflineHandler(msgID uint64, data []byte) {
	msg := &gateway.ClientOffline{}
	err := c.LogicRouter.Unmarshal(data, msg)
	if err != nil {
		logger.Error("receive data:%v msg:%v", data, msg)
	}
	pData := &player.PlayerData{
		PlayerID:  msg.Userid,
		GatewayID: c.serverID,
		OpType:    msg.OpType,
	}
	logger.Info("[clientOfflineHandler] player offline UserId:%v下线消息 通知 world goroutine", msg.Userid)
	server.Oasis.ChanPlayerOffline <- pData
}
