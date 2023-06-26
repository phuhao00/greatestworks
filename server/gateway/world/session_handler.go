package world

import (
	"github.com/phuhao00/greatestworks-proto/chat"
	"github.com/phuhao00/greatestworks-proto/gateway"
	"github.com/phuhao00/greatestworks-proto/messageId"
	"github.com/phuhao00/greatestworks-proto/server_common"
	"greatestworks/aop/logger"
	"greatestworks/server/gateway/client"
	"greatestworks/server/gateway/server"
	"reflect"
)

func (s *Session) Register() {
	s.LogicRouter.Register(uint64(messageId.MessageId_RegisterToGateway), s.registerHandler)
	s.LogicRouter.Register(uint64(messageId.MessageId_UpdateOnlineInfo), s.updateHandler)
	s.LogicRouter.Register(uint64(messageId.MessageId_BroadcastChatMsg), s.broadcastChatMsgHandler)
	s.LogicRouter.Register(uint64(messageId.MessageId_ClientOnlineRet), s.clientOnlineHandler)
	s.LogicRouter.Register(uint64(messageId.MessageId_ClientOfflineRet), s.clientOfflineHandler)
	s.LogicRouter.Register(uint64(messageId.MessageId_GatewayForwardPacket), s.gatewayForwardPacketHandler)
}

func (s *Session) registerHandler(msgID uint64, data []byte) {
	msg := &server_common.RegisterToGateway{}
	err := s.Router.Unmarshal(data, msg)
	if err != nil {
		logger.Error("[registerHandler] receive data:%v msg:%v", data, msg)
		return
	}
	if msg.Token != "greatest_works" {
		return
	}
	localZoneId := server.GetServer().Config.Global.ZoneId
	if int(msg.ZoneId) != localZoneId {
		logger.Warn("[registerHandler] ServerType:%v ServerAddr:%v, 但是zone 本地 %v != 远端 %v ",
			msg.ServerType, msg.ServerAddr, msg.ProcIndex, localZoneId, msg.ZoneId)
		return
	}
	logger.Info("[registerHandler] ServerType:%v ServerAddr:%v proIndex:%v, zone: %v",
		msg.ServerType, msg.ServerAddr, msg.ProcIndex, msg.ZoneId)

	switch msg.ServerType {
	case server_common.ServerCategory_World:
		GetMe().AddSession(msg.ServerID, msg.ServerAddr, msg.MsgIds, s, msg.ProcIndex, int(msg.ZoneId))
	default:
		logger.Warn("[registerHandler] ServerType:%v ServerAddr:%v, zone: %v",
			msg.ServerType, msg.ServerAddr, msg.ZoneId)
	}
	s.Name = msg.Name
	s.MaxPlayer = msg.MaxPlayer
	s.Verify()
	msgSend := &server_common.RegisterToGatewayRet{Result: uint32(server_common.ServerRegisterStatus_OK), ZoneId: msg.ZoneId}
	s.AsyncSend(uint64(messageId.MessageId_RegisterToGatewayRet), msgSend)
}

func (s *Session) updateHandler(msgID uint64, data []byte) {
	msg := &server_common.UpdateOnlineInfo{}
	err := s.Router.Unmarshal(data, msg)
	if err != nil {
		logger.Error("[updateHandler] receive data:%v msg:%v", data, msg)
		return
	}
	s.MaxPlayer = msg.MaxPlayer
	msgSend := &server_common.UpdateOnlineInfoRet{Result: uint32(server_common.ServerRegisterStatus_OK)}
	s.AsyncSend(uint64(messageId.MessageId_RegisterToGatewayRet), msgSend)
	logger.Info("[updateHandler]  name:%v ,MaxPlayer:%v", s.Name, s.MaxPlayer)
}

// broadcastChatMsgHandler 广播聊天消息
func (s *Session) broadcastChatMsgHandler(msgID uint64, data []byte) {
	msg := &chat.SCCrossSrvChatMsg{}
	err := s.Router.Unmarshal(data, msg)
	if err != nil {
		logger.Error("[broadcastChatMsgHandler] receive data:%v msg:%v", data, msg)
		return
	}
	if msg.MessageType == 2 {
		client.GetMe().BroadcastMsgSrv(messageId.MessageId_SCCrossSrvChatMsg, msg)
	} else {
		client.GetMe().BroadcastMsg(messageId.MessageId_SCCrossSrvChatMsg, msg)
	}
}

// clientOnlineHandler 玩家上线 处理
func (s *Session) clientOnlineHandler(msgID uint64, data []byte) {
	msg := &gateway.ClientOnlineRet{}
	err := s.Router.Unmarshal(data, msg)
	if err != nil {
		logger.Error("[clientOnlineHandler] receive data:%v msg:%v", data, reflect.TypeOf(msg))
		return
	}
	s.Players = msg.GetPlayers()
	client.GetMe().ClientOnline(msg, s.SessionId)
}

// clientOfflineHandler 玩家离线 处理
func (s *Session) clientOfflineHandler(msgID uint64, data []byte) {
	msg := &gateway.ClientOfflineRet{}
	err := s.Router.Unmarshal(data, msg)
	if err != nil {
		logger.Error("[clientOnlineHandler] receive data:%v, msg:%v", data, reflect.TypeOf(msg))
		return
	}
	if msg.Result != 1 {
		logger.Error("[clientOnlineHandler] offline error result:%v ,clientID:%v", msg.Result, msg.Userid)
		return
	}
	s.Players = msg.GetPlayers()
	clientInstance := client.GetMe().GetClientByUserID(msg.Userid)
	if msg.OpType == 1 {
		if clientInstance != nil {
			clientInstance.ClientLogout()
		}
		logger.Info("[clientOnlineHandler]  userID:%v ,address1:%v ,online:%v", msg.Userid, s.serverAddr, msg.ProcIndex)
	} else if clientInstance != nil {
		joinOnlineID := clientInstance.JoinWorldServerId.Load().(string)
		clientInstance.SendClientOnlineMsg(joinOnlineID, msg.Version, false)
	}
}

// gatewayForwardPacketHandler  world server 的消息发给客户端
func (s *Session) gatewayForwardPacketHandler(msgID uint64, data []byte) {
	msg := &gateway.GatewayForwardPacket{}
	err := s.Router.Unmarshal(data, msg)
	if err != nil {
		logger.Error("[gatewayForwardPacketHandler] receive data:%v msg:%v", data, msg)
		return
	}
	clientInstance := client.GetMe().GetClientByUserID(msg.Userid)
	if clientInstance != nil && clientInstance.IsBindWorldServer.Load().(bool) {
		sendData := make([]byte, len(msg.Data))
		copy(sendData, msg.Data)
		client.GetMe().SendMsgToClient(msg.Userid, sendData)

		if msg.PlayerNumber >= 0 {
			s.Players = msg.PlayerNumber
		}
	}
}
