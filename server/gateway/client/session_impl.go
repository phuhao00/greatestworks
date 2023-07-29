package client

import (
	"github.com/phuhao00/network"
	"github.com/phuhao00/spoor/logger"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func (s *Session) OnConnect() {
	GetMe().addClient(s.ConnID, s)
	logger.Info("[OnConnect]  local:%s remote:%s ConnID:%v", s.LocalAddr(), s.RemoteAddr(), s.ConnID)
	info := strings.Split(s.RemoteAddr().String(), ":")
	if len(info) > 0 {
		s.RemoteIp = info[0]
	} else {
		s.RemoteIp = s.RemoteAddr().String()
	}
}

func NewSession(session *network.TcpSession) network.ISession {
	return &Session{
		TcpSession:           session,
		Router:               nil,
		LogicRouter:          nil,
		UserID:               0,
		CharacterId:          0,
		WorldServerId:        atomic.Value{},
		IsBindWorldServer:    atomic.Value{},
		JoinWorldServerId:    atomic.Value{},
		IsDisconnected:       atomic.Value{},
		DisconnectedTime:     atomic.Value{},
		DisconnectedMessages: nil,
		IsReconnection:       atomic.Value{},
		RemoteIp:             "",
		LastPingTime:         time.Time{},
		LastCheckTime:        0,
		RequestSpeed:         0,
		ReqFrequency:         0,
		MsgRegisterOnce:      sync.Once{},
		ProcIndex:            0,
		mu:                   sync.Mutex{},
	}
}
