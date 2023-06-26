package client

import (
	"context"
	"github.com/phuhao00/greatestworks-proto/chat"
	"github.com/phuhao00/greatestworks-proto/gateway"
	"github.com/phuhao00/greatestworks-proto/messageId"
	"greatestworks/aop/logger"
	"greatestworks/aop/redis"
	"greatestworks/server/gateway/server"
	"greatestworks/server/gateway/world"
	"sync"
	"sync/atomic"
	"time"
)

const (
	disconnectedDuration   = 60 * time.Second
	pingNoResponseDuration = 5 * time.Minute
)

var (
	initOnce      sync.Once
	clientManager *Manager
)

type Manager struct {
	clients        *sync.Map
	userid2Clients *sync.Map
	players        int32
}

func GetMe() *Manager {
	initOnce.Do(func() {
		clientManager = &Manager{
			clients:        &sync.Map{},
			userid2Clients: &sync.Map{},
		}
	})
	return clientManager
}

func (m *Manager) GetClientNum() uint32 {
	var count uint32
	count = 0
	m.clients.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

func (m *Manager) addClient(clientID int64, c *Session) {
	_, ok := m.clients.Load(clientID)
	if ok {
		logger.Warn("[addClient] clientID: %v already exist", clientID)
	}
	m.clients.Store(clientID, c)
	atomic.AddInt32(&m.players, 1)
}

func (m *Manager) removeClient(clientID int64) {
	_, ok := m.clients.Load(clientID)
	if !ok {
		return
	}
	m.clients.Delete(clientID)
	atomic.AddInt32(&m.players, -1)
}

func (m *Manager) bindUserID2Client(userID uint64, client *Session) {
	_, ok := m.userid2Clients.Load(userID)
	if ok {
		logger.Warn("[bindUserID2Client] userId: %v already exist", userID)
	}
	m.userid2Clients.Store(userID, client)
}

func (m *Manager) unbindUserID2Client(userID uint64, cli *Session) {
	c, ok := m.userid2Clients.Load(userID)
	if !ok {
		return
	}

	if c.(*Session) == cli {
		m.userid2Clients.Delete(userID)
	}
}

func (m *Manager) getClient(clientID int64) *Session {
	client, ok := m.clients.Load(clientID)
	if !ok {
		return nil
	}
	return client.(*Session)
}

func (m *Manager) GetClientByUserID(userID uint64) *Session {
	client, ok := m.userid2Clients.Load(userID)
	if !ok {
		return nil
	}
	return client.(*Session)
}

func (m *Manager) SendMsgToClient(userID uint64, data []byte) bool {
	client := m.GetClientByUserID(userID)
	if client == nil {
		return false
	}

	if !client.IsDisconnected.Load().(bool) {
		return client.AsyncSendRowMsg(data)
	}

	if time.Since(client.DisconnectedTime.Load().(time.Time)) < disconnectedDuration {
		client.mu.Lock()
		if int64(len(client.DisconnectedMessages)) < server.GetServer().PriMsgBuffSize {
			client.DisconnectedMessages = append(client.DisconnectedMessages, data)
		} else {
			client.clientDisConnection()
		}
		client.mu.Unlock()
		return false
	} else {
		client.clientDisConnection()
	}

	return false
}

// ClientOnline 把消息号id映射到指定的handler
func (m *Manager) ClientOnline(msg *gateway.ClientOnlineRet, srvID string) {
	client := m.GetClientByUserID(msg.Userid)
	if client == nil {
		logger.Error("[ClientOnline] can not found client userID:%v", msg.Userid)
		return
	}
	worldServer, err := world.GetMe().GetEndpoint(srvID)
	if err != nil {
		logger.Error("[ClientOnline] error:", err)
		return
	}

	client.MsgRegisterOnce.Do(func() {
		for _, msgID := range worldServer.HandleMessages {
			client.registerForwardMsgHandler(msgID, world.GetMe().WorldMsgHandler)
		}
	})

	if msg.Result == 1 {
		client.JoinWorldServerId.Store("")
		client.IsBindWorldServer.Store(true)
		client.WorldServerId.Store(worldServer.SessionId)
	}

	if msg.Reconnection {
		client.clientReconnection(msg.Result)
		client.IsReconnection.Store(false)
		return
	}

	msgSend := &gateway.SCGatewayJoinOnline{
		Userid:          msg.Userid,
		Nick:            msg.Name,
		Gold:            msg.Gold,
		SceneId:         msg.SceneId,
		IsNew:           msg.IsNew,
		Frame:           msg.Frame,
		Head:            msg.Head,
		Model:           msg.Model,
		Sex:             msg.Sex,
		Level:           msg.Level,
		Exp:             msg.Exp,
		RegTime:         msg.RegTime,
		ServerTime:      msg.ServerTime,
		TodayFirstLogin: msg.TodayFirstLogin,
		ProcIndex:       msg.ProcIndex,
		ErrCode:         gateway.JoinErrorCode_JoinWorldSuccess,
	}
	client.sendMsg(messageId.MessageId_SCGatewayJoinWorld, msgSend)
	client.SetProcIndex(msg.ProcIndex)
	logger.Info("[ClientOnline] userID:%v isNew %v online:%v, clientIP: %v", msg.Userid, msg.IsNew, msg.ProcIndex, client.RemoteIp)
}

func (m *Manager) getPlayers() int32 {
	return m.players
}

func (m *Manager) WorldServerDisconnected(srvID, srvAddr string, zoneId int, proIndx uint32) {
	m.clients.Range(func(k, v interface{}) bool {
		client, ok := v.(*Session)
		if ok && client.WorldServerId.Load().(string) == srvID {
			if client.UserID != 0 {
				logger.Info("[WorldServerDisconnected] srvAddr:%v zoneId:%v pidx:%v  uid:%v", srvAddr, zoneId, proIndx, client.UserID)
			}
			GetMe().removeClient(client.ConnID)
			GetMe().unbindUserID2Client(client.UserID, client)
			client.UserID = 0
			client.Close()
		}
		return true
	})
}

func (m *Manager) DisconnectedCheck() {
	m.clients.Range(func(k, v interface{}) bool {
		client, ok := v.(*Session)
		if ok {
			if client.IsDisconnected.Load().(bool) &&
				time.Since(client.DisconnectedTime.Load().(time.Time)) > disconnectedDuration {
				client.clientDisConnection()
			}

			if time.Since(client.LastPingTime) > pingNoResponseDuration {
				client.clientNoResponseTimeout()
			}
		}
		return true
	})
}

func (m *Manager) SendMsgToAllPlayer(cmd messageId.MessageId, msg interface{}) {
	m.clients.Range(func(k, v interface{}) bool {
		client, ok := v.(*Session)
		if ok {
			if cmd == messageId.MessageId_SCKick {
				client.sendLastMsg(cmd, msg)
				if client.WorldServerId.Load().(string) != "" {
					client.sendClientOfflineMsg(1, client.WorldServerId.Load().(string))
				}
				key := redis.MakeAccountKey(int64(client.UserID))
				ret := redis.CacheRedis().HGet(context.TODO(), key, "Gateway")
				if ret.Val() == server.GetServer().GetTcpAddress() {
					redis.CacheRedis().HDel(context.TODO(), key, "Gateway", "ZoneId")
					logger.Debug("[SendMsgToAllPlayer] redis HDel ZoneId of %v", key)
				}
				GetMe().removeClient(client.ConnID)
				GetMe().unbindUserID2Client(client.UserID, client)
			} else {
				client.sendMsg(cmd, msg)
			}
		}
		return true
	})
}

func (m *Manager) BroadcastMsg(cmd messageId.MessageId, msg interface{}) {
	m.clients.Range(func(k, v interface{}) bool {
		client, ok := v.(*Session)
		if ok {
			client.sendMsg(cmd, msg)
		}
		return true
	})
}

func (m *Manager) BroadcastMsgSrv(cmd messageId.MessageId, msg interface{}) {
	message, ok := msg.(*chat.SCCrossSrvChatMsg)
	if !ok {
		logger.Error("[BroadcastMsgSrv]  assert fail ，not SCCrossSrvChatMsg , msg: %v", message)
		return
	}
	m.clients.Range(func(k, v interface{}) bool {
		client, ok := v.(*Session)
		if ok && (client.ProcIndex/message.RangeOfSrv == message.ProcIndex/message.RangeOfSrv) {
			client.sendMsg(cmd, msg)
		}
		return true
	})
}
