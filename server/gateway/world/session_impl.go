package world

import (
	"github.com/phuhao00/network"
	"github.com/phuhao00/spoor/logger"
)

func (s *Session) OnConnect() {
	logger.Info("[OnConnect]  local:%s remote:%s ConnID:%v", s.LocalAddr(), s.RemoteAddr(), s.ConnID)
}

func NewSession(session *network.TcpSession) network.ISession {
	return &Session{
		TcpSession:     session,
		Router:         nil,
		LogicRouter:    nil,
		serverType:     0,
		SessionId:      "",
		serverAddr:     "",
		HandleMessages: nil,
		Players:        0,
		ProIndex:       0,
		Name:           "",
		MaxPlayer:      0,
		ZoneId:         0,
	}

}
