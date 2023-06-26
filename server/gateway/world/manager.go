package world

import (
	"errors"
	"github.com/phuhao00/greatestworks-proto/gateway"
	"github.com/phuhao00/greatestworks-proto/messageId"
	"github.com/phuhao00/greatestworks-proto/server_common"
	"greatestworks/aop/logger"
	"greatestworks/server/gateway/client"
	"greatestworks/server/gateway/server"
	"math/rand"
	"sort"
	"sync"
)

var (
	onceInit           sync.Once
	worldServerManager *Manager
	ErrNoEndpoints     = errors.New("no endpoints available")
)

type Manager struct {
	Endpoints sync.Map
}

func GetMe() *Manager {
	onceInit.Do(func() {
		worldServerManager = &Manager{Endpoints: sync.Map{}}
	})
	return worldServerManager
}

// GetSessionList ...
func (m *Manager) GetSessionList() []*Session {
	endpoints := make([]*Session, 0, 128)
	m.Endpoints.Range(func(key, value interface{}) bool {
		serverClient, ok := value.(*Session)
		if ok {
			endpoints = append(endpoints, serverClient)
		}
		return true
	})
	if len(endpoints) > 0 {
		sort.Slice(endpoints[:], func(i, j int) bool {
			return endpoints[i].ProIndex < endpoints[j].ProIndex
		})
	}
	return endpoints
}

func (m *Manager) AddSession(srvID, address string, msgIds []uint32, sc *Session, proIndex uint32, zoneId int) {
	oldc, ok := m.Endpoints.Load(srvID)
	if ok {
		oldClient, _ := oldc.(*Session)
		if oldClient != sc {
			oldClient.Close()
		}
		logger.Error("[AddSession] 连接已经注册成功,重复注册 srvID:%v addr:%v", srvID, address)
	}
	sc.serverAddr = address
	sc.SessionId = srvID
	sc.serverType = server_common.ServerCategory_World
	sc.ProIndex = proIndex
	sc.ZoneId = zoneId
	for _, msgID := range msgIds {
		sc.HandleMessages = append(sc.HandleMessages, uint64(msgID))
	}
	m.Endpoints.Store(srvID, sc)
}

// RemoveServerClient ...
func (m *Manager) RemoveServerClient(srvID string) {
	s, ok := m.Endpoints.Load(srvID)
	if !ok {
		return
	}
	m.Endpoints.Delete(srvID)
	srv := s.(*Session)
	client.GetMe().WorldServerDisconnected(srvID, srv.serverAddr, srv.ZoneId, srv.ProIndex)
}

func (m *Manager) GetEndpoint(srvID string) (*Session, error) {
	session, ok := m.Endpoints.Load(srvID)
	if !ok {
		zoneId := server.GetServer().Config.Global.ZoneId
		logger.Warn("[GetEndpoint] try get scene server error unValuable service endpoint, srvID:%v in this zone %v", srvID, zoneId)
		return nil, ErrNoEndpoints
	}
	return session.(*Session), nil
}

func (m *Manager) SendMsgToWorldServer(srvID string, cmd messageId.MessageId, msg interface{}) bool {
	session, ok := m.Endpoints.Load(srvID)
	if !ok {
		return false
	}
	return session.(*Session).sendMsg(cmd, msg)

}

func (m *Manager) GetOptimalOnline() (string, bool) {
	idle := make([]*Session, 0)
	busy := make([]*Session, 0)
	hot := make([]*Session, 0)
	for _, online := range m.GetSessionList() {
		rate := (float32(online.Players) / float32(online.MaxPlayer)) * 100
		if rate < 100 && rate > 80 {
			hot = append(hot, online)
		} else if rate <= 80 && rate > 40 {
			busy = append(busy, online)
		} else if rate <= 40 {
			idle = append(idle, online)
		}
	}
	if len(busy) > 0 {
		return busy[rand.Int()%len(busy)].SessionId, true
	}
	if len(idle) > 0 {
		return idle[rand.Int()%len(idle)].SessionId, true
	}
	if len(hot) > 0 {
		return hot[rand.Int()%len(hot)].SessionId, true
	}
	return "", false
}

func (m *Manager) WorldMsgHandler(srvID string, userID uint64, msgID uint64, data []byte) {
	serverClient, ok := m.Endpoints.Load(srvID)
	if !ok {
		zoneId := server.GetServer().Config.Global.ZoneId
		logger.Warn("[WorldMsgHandler] unknown server may be server is offline, srvID:%v, local zone %v", srvID, zoneId)
		return
	}
	msgSend := &gateway.GatewayForwardPacket{Userid: userID}
	msgSend.Data = make([]byte, len(data))
	copy(msgSend.Data, data)
	serverClient.(*Session).sendMsg(messageId.MessageId_GatewayForwardPacket, msgSend)
}
