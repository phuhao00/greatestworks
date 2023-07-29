package client

import (
	"context"
	"github.com/phuhao00/fuse"
	"github.com/phuhao00/greatestworks-proto/gateway"
	"github.com/phuhao00/greatestworks-proto/messageId"
	"github.com/phuhao00/network"
	"greatestworks/aop/logger"
	"greatestworks/aop/redis"
	"greatestworks/server/gateway/server"
	"greatestworks/server/gateway/world"
	"sync"
	"sync/atomic"
	"time"
)

const (
	CheckIntervalSec   = 30       // 计算单个客户端请求流量的时间跨度 Sec
	ClientMaxSpeed     = 1024 * 5 // 单个客户端允许的最大请求流量 Byte/Sec
	ClientMaxFrequency = 30       // 单个客户端允许的最大请求频率 Times/Sec
)

type Session struct {
	*network.TcpSession
	Router               *fuse.Router      // 消息路由器
	LogicRouter          *fuse.LogicRouter // 消息路由器
	UserID               uint64            // uint64
	CharacterId          uint64            // 角色Id
	WorldServerId        atomic.Value      // string
	IsBindWorldServer    atomic.Value      // bool
	JoinWorldServerId    atomic.Value      // string
	IsDisconnected       atomic.Value      // bool
	DisconnectedTime     atomic.Value      // time.Time
	DisconnectedMessages []interface{}     // 断线缓存消息
	IsReconnection       atomic.Value      // bool
	RemoteIp             string            // 对端IP
	LastPingTime         time.Time         // 上次ping的时间
	LastCheckTime        int64             // 上次收到数据的时间点
	RequestSpeed         int               // 平均上行流速
	ReqFrequency         int               // 平均上行评率
	MsgRegisterOnce      sync.Once         // 代理消息注册一次
	ProcIndex            uint32            // 服id编号
	mu                   sync.Mutex        // mutex

}

func (s *Session) Init() {
	now := time.Now()
	s.Router = fuse.NewRouter()
	s.CharacterId = 0
	s.LastCheckTime = now.Unix()
	s.RequestSpeed = 0
	s.ReqFrequency = 0
	s.WorldServerId.Store("")
	s.JoinWorldServerId.Store("")
	s.IsBindWorldServer.Store(false)
	s.IsDisconnected.Store(false)
	s.IsReconnection.Store(false)
	s.DisconnectedTime.Store(now)
	s.LastPingTime = now
}

func (s *Session) Resolve(*network.Packet) {

}

func (s *Session) Loop() {

	for {
		select {}
	}
}

func (s *Session) registerForwardMsgHandler(msgID uint64, handler fuse.ForwardMessageHandler) {
	s.LogicRouter.RegisterForwardHandler(msgID, handler)
}

func (s *Session) HandleMessage(data []byte) {
	now := time.Now().Unix()
	delta := int(now - s.LastCheckTime)

	s.RequestSpeed += len(data)
	s.ReqFrequency++

	if delta >= CheckIntervalSec {
		speed := s.RequestSpeed / delta
		freq := s.ReqFrequency / delta

		s.LastCheckTime = now
		s.RequestSpeed = 0
		s.ReqFrequency = 0

		if speed >= ClientMaxSpeed || freq >= ClientMaxFrequency {
			logger.Error("[OnMessage]客户端疑似外挂, 请求流量:%v(字节/秒), 频率:%v(次/秒), [ip:%v, userid:%v, onlineId:%v]",
				speed, freq, s.RemoteIp, s.UserID, s.WorldServerId.Load())
			s.Close()
			return
		}

		if speed > 100 || freq > 5 {
			logger.Info("[OnMessage] 客户端请求流量:%v(字节/秒), 频率:%v(次/秒), [ip:%v, userid:%v, onlineId:%v]",
				speed, freq, s.RemoteIp, s.UserID, s.WorldServerId.Load())
		}
	}

	msgID, err := s.LogicRouter.Route(data)
	if err != nil && s.Verified() {
		msgID, err = s.LogicRouter.ForwardRoute(s.WorldServerId.Load().(string), s.UserID, data)
	}

	if err != nil {
		logger.Error("[OnMessage] 消息:%v路由失败 未注册该消息处理器 error: %v", msgID, err)
		return
	}
	if messageId.MessageId(msgID) != messageId.MessageId_SceneHeartbeat &&
		messageId.MessageId(msgID) != messageId.MessageId_CSPlayerMove {
		logger.Debug("[OnMessage] userId:%v 消息ID::%v", s.UserID, messageId.MessageId(msgID))
	}

	s.LastPingTime = time.Now()
}

// OnClose ...
func (s *Session) OnClose() {
	if s.UserID == 0 {
		GetMe().removeClient(s.ConnID)
		GetMe().unbindUserID2Client(s.UserID, s)
	} else {
		s.IsDisconnected.Store(true)
		s.DisconnectedTime.Store(time.Now())
	}
	logger.Info("[OnClose]  local:%v remote:%v userID:%v", s.LocalAddr(), s.RemoteAddr(), s.UserID)
}

func (s *Session) Marshal(msgID uint16, msg interface{}) ([]byte, error) {
	return s.Router.Marshal(msgID, msg)
}

func (s *Session) sendMsg(cmd messageId.MessageId, msg interface{}) bool {
	return s.AsyncSend(uint64(cmd), msg)
}

func (s *Session) sendLastMsg(cmd messageId.MessageId, msg interface{}) bool {
	return s.AsyncSendLastPacket(uint64(cmd), msg)
}

func (s *Session) onUserVerify(userID uint64) {
	s.UserID = userID
	GetMe().bindUserID2Client(userID, s)
	s.Verify()
}

func (s *Session) SendClientOnlineMsg(onlineID, version string, reconnection bool) {
	sendMsg := &gateway.ClientWorld{
		Userid:        s.UserID,
		IsReconnected: reconnection,
		RemoteIp:      s.RemoteIp,
		GatewayIp:     server.GetServer().GetTcpAddress(),
		Version:       version,
	}
	world.GetMe().SendMsgToWorldServer(onlineID, messageId.MessageId_ClientOnline, sendMsg)
	onlineServer, err := world.GetMe().GetEndpoint(onlineID)
	if err != nil {
		logger.Info("[SendClientOnlineMsg] UserID:%v, clientID:%v reConn:%v rip:%v online:%v", s.UserID, s.ConnID, reconnection, s.RemoteIp, onlineID)
	} else {
		logger.Info("[SendClientOnlineMsg] UserID:%v,clientID:%v reConn:%v rip:%v online:%v", s.UserID, s.ConnID, reconnection, s.RemoteIp, onlineServer.ProIndex)
	}
}

func (s *Session) sendClientOfflineMsg(opType uint32, worldServerId string) {
	sendMsg := &gateway.ClientOffline{
		Userid: s.UserID,
		OpType: opType,
	}
	world.GetMe().SendMsgToWorldServer(worldServerId, messageId.MessageId_ClientOffline, sendMsg)
	worldServer, err := world.GetMe().GetEndpoint(worldServerId)
	if err != nil {
		logger.Info("[sendClientOfflineMsg] UserID:%v, clientID:%v 下线类型:%v online:%v", s.UserID, s.ConnID, opType, worldServerId)
	} else {
		logger.Info("[sendClientOfflineMsg] UserID:%v, clientID:%v 下线类型:%v online:%v", s.UserID, s.ConnID, opType, worldServer.ProIndex)
	}
}

func (s *Session) getDisconnectedMessages() []interface{} {
	var messages []interface{}
	s.mu.Lock()
	messages = s.DisconnectedMessages
	s.mu.Unlock()
	return messages
}

func (s *Session) clientReconnection(ret uint32) {
	sendMsg := &gateway.SCReconnection{Ret: ret}
	s.sendMsg(messageId.MessageId_CSReconnection, sendMsg)

	var messagesLen int
	s.mu.Lock()
	messagesLen = len(s.DisconnectedMessages)
	for _, data := range s.DisconnectedMessages {
		bytesData, ok := data.([]byte)
		if !ok {
			continue
		}
		s.AsyncSendRowMsg(bytesData)
	}
	s.DisconnectedMessages = s.DisconnectedMessages[0:0]
	s.mu.Unlock()
	logger.Info("[clientReconnection] ret:%v 离线后缓存的消息数量:%v", ret, messagesLen)
}

func (s *Session) clientDisConnection() {
	if s.WorldServerId.Load().(string) != "" {
		s.sendClientOfflineMsg(1, s.WorldServerId.Load().(string))
	}

	key := redis.MakeAccountKey(int64(s.UserID))
	ret := redis.CacheRedis().HGet(context.TODO(), key, "Gateway")
	if ret.Val() == server.GetServer().GetTcpAddress() {
		redis.CacheRedis().HDel(context.TODO(), key, "Gateway", "ZoneId")
		logger.Debug("[clientDisConnection]redis HDel ZoneId of %v", key)
	}

	GetMe().removeClient(s.ConnID)
	GetMe().unbindUserID2Client(s.UserID, s)
	if !s.IsClosed() {
		s.Close()
	}
	redis.CacheRedis().SRem(context.TODO(), redis.MakeGatewayKey(server.GetServer().GetTcpAddress()), s.UserID)
	logger.Info("[clientDisConnection]  connID:%v userID:%v", s.ConnID, s.UserID)
}

func (s *Session) clientNoResponseTimeout() {
	if s.WorldServerId.Load().(string) != "" {
		s.sendClientOfflineMsg(1, s.WorldServerId.Load().(string))
	}
	key := redis.MakeAccountKey(int64(s.UserID))
	ret := redis.CacheRedis().HGet(context.TODO(), key, "Gateway")
	if ret.Val() == server.GetServer().GetTcpAddress() {
		redis.CacheRedis().HDel(context.TODO(), key, "Gateway", "ZoneId")
		logger.Debug("[clientNoResponseTimeout]redis HDel ZoneId of %v", key)
	}
	GetMe().removeClient(s.ConnID)
	GetMe().unbindUserID2Client(s.UserID, s)
	// 连接未断开的断开连接
	if !s.IsClosed() {
		s.Close()
	}
	redis.CacheRedis().SRem(context.TODO(), redis.MakeGatewayKey(server.GetServer().GetTcpAddress()), s.UserID)
	logger.Info("[clientNoResponseTimeout]   connID:%v userID:%v", s.ConnID, s.UserID)
}

func (s *Session) ClientLogout() {
	redis.CacheRedis().SRem(context.TODO(), redis.MakeGatewayKey(server.GetServer().GetTcpAddress()), s.UserID)
	key := redis.MakeAccountKey(int64(s.UserID))
	ret := redis.CacheRedis().HGet(context.TODO(), key, "Gateway")
	if ret.Val() == server.GetServer().GetTcpAddress() {
		redis.CacheRedis().HDel(context.TODO(), key, "Gateway", "ZoneId")
		logger.Debug("[ClientLogout]redis HDel ZoneId of %v", key)
	}

	GetMe().removeClient(s.ConnID)
	GetMe().unbindUserID2Client(s.UserID, s)
	msgSend := &gateway.SCGatewayLogout{}
	s.sendLastMsg(messageId.MessageId_CSGatewayLogout, msgSend)
}

func (s *Session) SetProcIndex(procIndex uint32) {
	s.ProcIndex = procIndex
}
