package client

import (
	"github.com/phuhao00/greatestworks-proto/server_common"
	"greatestworks/aop/logger"
	"greatestworks/server/world/server"
)

func (c *Client) registerToGameHandler(msgID uint64, data []byte) {
	msg := &server_common.RegisterToGameRet{}
	err := c.LogicRouter.Unmarshal(data, msg)
	if err != nil {
		logger.Error("[registerToGameHandler] receive data:%v msg:%v", data, msg)
		return
	}

	if msg.Result == uint32(server_common.ServerRegisterStatus_OK) {
		for _, msgID := range msg.MsgIds {
			c.handleMessages = append(c.handleMessages, uint16(msgID))
		}
		logger.Info("[registerToGameHandler] serverID:%v Addr:%v", c.serverID, c.serverAddr)
	}

}

func (c *Client) createGameHandler(msgID uint64, data []byte) {
	msg := &server_common.CreateGameRet{}
	err := c.LogicRouter.Unmarshal(data, msg)
	if err != nil {
		logger.Error("[createGameHandler] receive data:%v msg:%v", data, msg)
		return
	}
	player := server.Oasis.GetPlayer(msg.PlayerId)
	if player != nil {
		err = server.Oasis.HandleServerMsgPacket(c.serverID, msg.PlayerId, data)
		if err != nil {
			logger.Error("[createGameHandler] err:%v", err.Error())
		}
	}
}

func (c *Client) gameMsgToOnlinePlayerHandler(msgID uint64, data []byte) {
	msg := &server_common.SrvMsgToWorldPlayer{}
	err := c.LogicRouter.Unmarshal(data, msg)
	if err != nil {
		logger.Error("[gameMsgToOnlinePlayerHandler] receive data:%v msg:%v", data, msg)
		return
	}
	player := server.Oasis.GetPlayer(msg.PlayerId)
	if player != nil {
		err = server.Oasis.HandleServerMsgPacket(c.serverID, msg.PlayerId, msg.Data)
		if err != nil {
			logger.Error("[gameMsgToOnlinePlayerHandler] err:%v", err.Error())
		}
	}
}
