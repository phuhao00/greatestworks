package world

import (
	"github.com/phuhao00/fuse"
	"github.com/phuhao00/greatestworks-proto/messageId"
	"github.com/phuhao00/greatestworks-proto/server_common"
	"github.com/phuhao00/network"
	"greatestworks/aop/logger"
	"runtime/debug"
)

type Session struct {
	*network.TcpSession
	Router         *fuse.Router
	LogicRouter    *fuse.LogicRouter
	serverType     server_common.ServerCategory
	SessionId      string
	serverAddr     string
	HandleMessages []uint64
	Players        int32
	ProIndex       uint32
	Name           string
	MaxPlayer      int32
	ZoneId         int
}

func (s *Session) sendMsg(cmd messageId.MessageId, msg interface{}) bool {
	return s.AsyncSend(uint64(cmd), msg)
}

func (s *Session) Marshal(msgID uint16, msg interface{}) ([]byte, error) {
	return s.Router.Marshal(msgID, msg)
}

func (s *Session) OnConnect() {
	logger.Info("[OnConnect]  local:%s remote:%s ConnID:%v", s.LocalAddr(), s.RemoteAddr(), s.ConnID)
}

func (s *Session) OnMessage(data []byte) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("[OnMessage] panic ", err, "\n", string(debug.Stack()))
		}
	}()

	_, err := s.LogicRouter.Route(data)

	if err != nil {
		logger.Error("[OnMessage] route message error: %v", err)
	}
}

func (s *Session) OnClose() {
	if s.SessionId != "" {
		GetMe().RemoveServerClient(s.SessionId)
	}
	logger.Info("[OnClose] local:%v remote:%v", s.LocalAddr(), s.RemoteAddr())
}
