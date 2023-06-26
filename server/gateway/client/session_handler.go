package client

import (
	"context"
	"github.com/phuhao00/fuse"
	"github.com/phuhao00/greatestworks-proto/ErrCode"
	"github.com/phuhao00/greatestworks-proto/gateway"
	"github.com/phuhao00/greatestworks-proto/messageId"
	"github.com/phuhao00/greatestworks-proto/scene"
	"github.com/phuhao00/greatestworks-proto/server_common"
	"github.com/phuhao00/network"
	"github.com/phuhao00/spoor/logger"
	"google.golang.org/protobuf/proto"
	"greatestworks/aop/redis"
	"greatestworks/server/gateway/gm"
	"greatestworks/server/gateway/server"
	"greatestworks/server/gateway/world"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

func (s *Session) HandlerRegister() {
	s.Router.AddRoute(uint64(messageId.MessageId_CSGatewayLogin), loginHandler)
	s.Router.AddRoute(uint64(messageId.MessageId_CSGatewayLogout), logoutHandler)
	s.Router.AddRoute(uint64(messageId.MessageId_CSGatewayWorldList), onlineServerListHandler)
	s.Router.AddRoute(uint64(messageId.MessageId_CSGatewayJoinWorld), joinOnlineHandler)
	s.Router.AddRoute(uint64(messageId.MessageId_CSReconnection), reconnectionHandler)
	s.Router.AddRoute(uint64(messageId.MessageId_SceneHeartbeat), sceneHeartbeatHandler)
}

func loginHandler(packet *network.Packet, principal fuse.Principal) {
	session := principal.(*Session)
	if session.Verified() {
		logger.Error("[loginHandler] Verified failed %v %v", session.LocalAddr(), session.RemoteAddr())
		return
	}

	msg := &gateway.CSGatewayLogin{}
	err := proto.Unmarshal(packet.Msg.Data, msg)
	if err != nil {
		logger.Error("[loginHandler] receive data Unmarshal error:%v msg:%v, localAddr:%v, remoteAddr:%v.",
			err.Error(), msg, session.LocalAddr(), session.RemoteAddr())
		return
	}

	msgSend := &gateway.SCGatewayLogin{
		ErrCode: uint32(gateway.GatewayErr_Success),
		Userid:  msg.Userid,
		Name:    msg.Token,
		ZoneId:  msg.ZoneId,
	}
	localZoneId := server.GetServer().Config.Global.ZoneId
	if int(msg.ZoneId) != localZoneId {
		logger.Error("[loginHandler]  msg:%v %v %v, but zone local %v != req %v ",
			msg, session.LocalAddr(), session.RemoteAddr(), localZoneId, msg.ZoneId)
		msgSend.ErrCode = uint32(gateway.GatewayErr_ZoneIdError)
		session.sendMsg(messageId.MessageId_SCGatewayLogin, msgSend)
		return
	}

	key := redis.MakeTokenKey(msg.Userid)
	ret := redis.CacheRedis().Get(context.TODO(), key)
	if msg.Token == "" || ret.Val() == "" || ret.Val() != msg.Token {
		msgSend.ErrCode = uint32(ErrCode.ErrCode_gateway_verify)
		logger.Error("[loginHandler] err token:%v userID:%v token:%v", ret.Val(), msg.Userid, msg.Token)
		session.sendMsg(messageId.MessageId_SCGatewayLogin, msgSend)
		return
	}

	client := GetMe().GetClientByUserID(msg.Userid)
	if client != nil && session == client { // 重复登录验证
		msgSend.ErrCode = uint32(ErrCode.ErrCode_gateway_repeated_verify)
		logger.Error("[loginHandler] duplicate verification client:%v %v", client.ConnID, msg.Userid)
		session.sendMsg(messageId.MessageId_SCGatewayLogin, msgSend)
		return
	}

	accountKey := redis.MakeAccountKey(int64(msg.Userid))
	gatewayInfo := redis.CacheRedis().HGet(context.TODO(), accountKey, "Gateway")

	if gatewayInfo.Err() != nil || gatewayInfo.Val() != server.GetServer().GetTcpAddress() {
		remoteIp := strings.Split(session.Conn.RemoteAddr().String(), ":")[0]
		msgSend.ErrCode = uint32(ErrCode.ErrCode_gateway_verify)
		logger.Error("[loginHandler] 客户端疑似外挂, 检验gateway失败, userid account key: [%v], clientIp:%v, recommend gateway:[%v] but this:[%v]",
			accountKey, remoteIp, gatewayInfo.Val(), server.GetServer().GetTcpAddress())
		session.sendMsg(messageId.MessageId_SCGatewayLogin, msgSend)
		return
	}

	if client != nil && client != session {
		GetMe().unbindUserID2Client(client.UserID, client)
		onlineID := client.WorldServerId.Load().(string)
		if onlineID != "" {
			session.WorldServerId.Store(onlineID)
			atomic.StoreUint64(&client.UserID, 0)
			client.WorldServerId.Store("")
		}
		kickoutMsg := &server_common.SCKick{Kick: server_common.KickReason_RemoteLogin}
		client.sendLastMsg(messageId.MessageId_SCKick, kickoutMsg)
		GetMe().removeClient(client.ConnID)
	}
	session.onUserVerify(msg.Userid)
	msgSend.ErrCode = uint32(gateway.GatewayErr_Success)
	session.sendMsg(messageId.MessageId_SCGatewayLogin, msgSend)
	redis.CacheRedis().SAdd(context.TODO(), redis.MakeGatewayKey(server.GetServer().GetTcpAddress()), msg.Userid)
	logger.Info("[login] gateway登录验证成功 %v", msg.Userid)
}

func onlineServerListHandler(packet *network.Packet, principal fuse.Principal) {
	session := principal.(*Session)
	if !session.Verified() {
		return
	}
	msg := &gateway.CSGatewayOnlineList{}
	err := proto.Unmarshal(packet.Msg.Data, msg)
	if err != nil {

	}

	localZoneId := server.GetServer().Config.Global.ZoneId
	if msg.ZoneId > 0 && int(msg.ZoneId) != localZoneId {
		logger.Error("[onlineServerListHandler] msg:%v %v %v, but zone local %v != req %v ",
			msg, session.LocalAddr(), session.RemoteAddr(), localZoneId, msg.ZoneId)
		return
	}

	msgSend := &gateway.SCGatewayOnlineList{}
	for _, onlineServer := range world.GetMe().GetSessionList() {
		msgSend.List = append(msgSend.List, &gateway.OnlineList{
			Sid:     onlineServer.SessionId,
			Players: onlineServer.Players,
			ProId:   onlineServer.ProIndex,
			Name:    onlineServer.Name,
			Max:     onlineServer.MaxPlayer,
			ZoneId:  int32(onlineServer.ZoneId),
		})
	}
	logger.Debug("[onlineServerListHandler] 拉取服务器列表 userid %v len %v", session.UserID, len(msgSend.List))
	session.sendMsg(messageId.MessageId_SCGatewayWorldList, msgSend)
}

func joinOnlineHandler(packet *network.Packet, principal fuse.Principal) {
	session := principal.(*Session)
	if !session.Verified() {
		return
	}

	msg := &gateway.CSGatewayJoinOnline{}
	err := proto.Unmarshal(packet.Msg.Data, msg)
	if err != nil {
		logger.Error("[joinOnlineHandler] receive data:%v msg:%v %v %v ", packet.Msg.Data, msg, session.LocalAddr(), session.RemoteAddr())
		return
	}

	logger.Info("[joinOnlineHandler] 玩家:%v请求进入在线服务器:%v 快速进入:%v", session.UserID, msg.Sid, msg.Quick)
	if msg.Quick {
		if sid, ok := world.GetMe().GetOptimalOnline(); ok {
			msg.Sid = sid
		}
	}

	msgSend := &gateway.SCGatewayJoinOnline{Userid: session.UserID}
	if !gm.GMInstance.IsOpenNow {
		remoteIp := strings.Split(session.Conn.RemoteAddr().String(), ":")[0]
		strUid := strconv.FormatInt(int64(session.UserID), 10)
		if !gm.GMInstance.IsIpInWhiteList(remoteIp) && !gm.GMInstance.IsUidInWhiteList(strUid) {
			logger.Info("[joinOnlineHandler] Now door closed, but client: %v ConnID: %v, uid:%v not in whitelist", remoteIp, session.ConnID, session.UserID)
			msgSend.ErrCode = gateway.JoinErrorCode_JoinWorldNotOpen
			session.sendMsg(messageId.MessageId_SCGatewayJoinWorld, msgSend)
			return
		}
		logger.Info("[joinOnlineHandler] Now door closed, and client ip:%v userid:%v in whitelist", remoteIp, session.UserID)
	}

	localZoneId := server.GetServer().Config.Global.ZoneId
	if int(msg.ZoneId) != localZoneId {
		logger.Error("[joinOnlineHandler] msg:%v %v %v, but zone local %v != req %v ",
			msg, session.LocalAddr(), session.RemoteAddr(), localZoneId, msg.ZoneId)
		// 客户端修改后加上检查条件
		msgSend.ErrCode = gateway.JoinErrorCode_JoinZoneIdError
		session.sendMsg(messageId.MessageId_SCGatewayLogin, msgSend)
		return
	}

	onlineServer, err := world.GetMe().GetEndpoint(msg.Sid)
	if err != nil {
		msgSend.ErrCode = gateway.JoinErrorCode_JoinWorldNotOpen
		session.sendMsg(messageId.MessageId_SCGatewayJoinWorld, msgSend)
		logger.Error("玩家:%v请求进去的在线:%v服错误:%v", session.UserID, msg.Sid, err)
		return
	}

	if session.JoinWorldServerId.Load().(string) != "" {
		msgSend.ErrCode = gateway.JoinErrorCode_JoinWorldIng
		session.sendMsg(messageId.MessageId_SCGatewayJoinWorld, msgSend)
		logger.Error("[joinOnlineHandler] 玩家:%v正处于上线状态 等待上线完成", session.UserID)
		return
	}

	if onlineServer.Players >= onlineServer.MaxPlayer {
		msgSend.ErrCode = gateway.JoinErrorCode_JoinWorldIsFull
		session.sendMsg(messageId.MessageId_SCGatewayJoinWorld, msgSend)
		logger.Error("[joinOnlineHandler] 玩家:%v请求进入服务器:%v 人数已满 最大人数:%v", session.UserID, msg.Sid, onlineServer.MaxPlayer)
		return
	}

	session.JoinWorldServerId.Store(onlineServer.SessionId)
	if session.WorldServerId.Load().(string) == "" {
		session.SendClientOnlineMsg(onlineServer.SessionId, msg.Version, false)
	} else {
		session.sendClientOfflineMsg(2, session.WorldServerId.Load().(string))
		session.WorldServerId.Store("")
	}
}

func reconnectionHandler(packet *network.Packet, principal fuse.Principal) {
	session := principal.(*Session)
	msg := &gateway.CSReconnection{}
	err := proto.Unmarshal(packet.Msg.Data, msg)
	if err != nil {
		logger.Error("[reconnectionHandler] receive data:%v msg:%v", packet.Msg.Data, msg)
		return
	}

	msgSend := &gateway.SCReconnection{Ret: uint32(gateway.GatewayErr_Success)}
	localZoneId := server.GetServer().Config.Global.ZoneId
	if int(msg.ZoneId) != localZoneId {
		logger.Error("[reconnectionHandler] msg:%v %v %v, but zone local %v != req %v ",
			msg, session.LocalAddr(), session.RemoteAddr(), localZoneId, msg.ZoneId)
		msgSend.Ret = uint32(gateway.JoinErrorCode_JoinZoneIdError)
		session.sendMsg(messageId.MessageId_SCReconnection, msgSend)
		return
	}

	logger.Info("[reconnectionHandler]玩家:%v重新连接 clientID:%v", msg.Userid, session.ConnID)
	if session.IsReconnection.Load().(bool) {
		logger.Debug("[reconnectionHandler] 玩家:%v正在重连中,请稍候...", msg.Userid)
		msgSend.Ret = uint32(gateway.GatewayErr_BeingReConn)
		session.sendMsg(messageId.MessageId_SCReconnection, msgSend)
		return
	}

	key := redis.MakeTokenKey(msg.Userid)
	ret := redis.CacheRedis().Get(context.TODO(), key)
	if msg.Token == "" || ret.Val() == "" || ret.Val() != msg.Token {
		msgSend.Ret = uint32(gateway.GatewayErr_Verify)
		session.sendMsg(messageId.MessageId_SCReconnection, msgSend)
		logger.Error("[reconnectionHandler] fail token:%v userID:%v token:%v", ret.Val(), msg.Userid, msg.Token)
		return
	}

	client := GetMe().GetClientByUserID(msg.Userid)
	if client == nil {
		msgSend.Ret = uint32(gateway.GatewayErr_LoseCurCon)
		session.sendMsg(messageId.MessageId_SCReconnection, msgSend)
		logger.Error("[reconnectionHandler] fail userID:%v", msg.Userid)
		return
	}

	if client != session {
		session.mu.Lock()
		session.DisconnectedMessages = client.getDisconnectedMessages()
		session.mu.Unlock()
		GetMe().removeClient(client.ConnID)
		GetMe().unbindUserID2Client(client.UserID, client)
		if !client.IsClosed() {
			client.UserID = 0
			client.Close()
		}
		session.onUserVerify(msg.Userid)
		session.SendClientOnlineMsg(client.WorldServerId.Load().(string), msg.Version, true)
		session.IsReconnection.Store(true) // 正在重连
		return
	}
	session.sendMsg(messageId.MessageId_SCReconnection, msgSend)
}

func sceneHeartbeatHandler(packet *network.Packet, principal fuse.Principal) {
	session := principal.(*Session)
	msg := &scene.SceneHeartbeat{}
	err := proto.Unmarshal(packet.Msg.Data, msg)
	if err != nil {
		logger.Error("[sceneHeartbeatHandler] receive data:%v msg:%v", packet.Msg.Data, msg)
		return
	}
	nowTime := time.Now()
	msg.Time = nowTime.Unix()
	session.sendMsg(messageId.MessageId_SceneHeartbeat, msg)
	session.LastPingTime = nowTime
}

func logoutHandler(packet *network.Packet, principal fuse.Principal) {
	session := principal.(*Session)
	if !session.Verified() {
		return
	}
	msg := &gateway.CSGatewayLogout{}
	err := proto.Unmarshal(packet.Msg.Data, msg)
	if err != nil {
		logger.Error("[logoutHandler] receive data:%v msg:%v %v %v ", packet.Msg.Data, msg, session.LocalAddr(), session.RemoteAddr())
		return
	}

	key := redis.MakeTokenKey(msg.UserId)
	ret := redis.CacheRedis().Get(context.TODO(), key)
	if msg.Token == "" || ret.Val() == "" || ret.Val() != msg.Token {
		logger.Error("[logoutHandler] 登录验证失败 token:%v userID:%v token:%v", ret.Val(), msg.UserId, msg.Token)
		return
	}

	if session.WorldServerId.Load().(string) != "" {
		session.sendClientOfflineMsg(1, session.WorldServerId.Load().(string))
	}
	redis.CacheRedis().SRem(context.TODO(), redis.MakeGatewayKey(server.GetServer().GetTcpAddress()), session.UserID)

	logger.Info("[logoutHandler] 退出游戏 %v", msg.UserId)
}
