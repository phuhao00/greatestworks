package server

import (
	"fmt"
	"github.com/phuhao00/broker/timerassistant"
	"github.com/phuhao00/fuse"
	pbChat "github.com/phuhao00/greatestworks-proto/chat"
	"github.com/phuhao00/greatestworks-proto/messageId"
	pbPLayer "github.com/phuhao00/greatestworks-proto/player"
	"github.com/phuhao00/greatestworks-proto/server_common"
	"greatestworks/aop/logger"
	"greatestworks/internal/communicate/chat"
	"greatestworks/internal/communicate/family"
	"greatestworks/internal/communicate/player"
	"greatestworks/server"
	"greatestworks/server/world/config"
	"os"
	"syscall"

	"github.com/phuhao00/network"
)

// World  todo 字节对齐
type World struct {
	*server.BaseService
	Pid               int
	Server            *network.TcpServer
	Handlers          map[messageId.MessageId]func(message *network.Packet)
	chSessionPacket   chan *network.Packet
	chatSystem        *chat.System
	familyManager     *family.Module
	playerManager     *player.Module
	timer             timerassistant.TimerAssistant
	LogicRouter       *fuse.LogicRouter
	ChanPlayerOffline chan *pbPLayer.PlayerData
	Config            *config.Config
	httpHandler       *HTTPHandler
	rpcAddr           string
	rpcPort           int
	httpAddr          string
	httpPort          int
	httpSvcID         string
	StartTM           int64
	crossZoneChatMsg  chan *pbChat.SCCrossSrvChatMsg //跨区聊天
	systemMsgChan     chan *pbChat.SCSystemMessage
}

func NewWorld() *World {
	m := &World{playerManager: player.NewPlayerMgr()}
	m.Server = network.NewTcpServer(":8023", 100, 200, logger.GetLogger())
	m.Server.MessageHandler = m.OnSessionPacket
	m.Handlers = make(map[messageId.MessageId]func(message *network.Packet))
	return m
}

var Oasis *World

func (w *World) OnSessionPacket(packet *network.Packet) {
	if handler, ok := w.Handlers[messageId.MessageId(packet.Msg.ID)]; ok {
		handler(packet)
		return
	}
	if p := w.playerManager.GetPlayer(uint64(packet.Conn.ConnID)); p != nil {
		p.HandlerParamCh <- packet.Msg
	}
}

func (w *World) OnSystemSignal(signal os.Signal) bool {
	logger.Debug("[OnSystemSignal] signal: %v \n", signal)
	tag := true
	switch signal {
	case syscall.SIGHUP:
		//todo
	case syscall.SIGPIPE:
	default:
		logger.Debug("[OnSystemSignal] ready exit...")
		tag = false

	}
	return tag
}

func (w *World) HandlePlayersMsgPacket(playerID uint64, data []byte) error {
	p := w.GetPlayer(playerID)
	if p != nil {
		msgData := &pbPLayer.PlayerMsgData{PlayerID: playerID}
		msgData.Data = make([]byte, len(data))
		copy(msgData.Data, data)
		p.HandleClientMsgPacket(msgData)
		return nil
	}
	return fmt.Errorf("[HandlePlayersMsgPacket] player not exist")
}

func (w *World) GetPlayer(playerID uint64) *player.Player {
	return w.playerManager.GetPlayer(playerID)
}

func (w *World) GetPlayersNum() int {
	return w.playerManager.GetPlayerNum()
}

func (w *World) HandleServerMsgPacket(srvID string, playerID uint64, data []byte) error {
	p := w.GetPlayer(playerID)
	if p != nil {
		serverMsgData := &server_common.ServerMsgData{ServerID: srvID, PlayerID: playerID}
		serverMsgData.Data = make([]byte, len(data))
		copy(serverMsgData.Data, data)
		p.HandleServerMsgPacket(serverMsgData)
		return nil
	}
	return fmt.Errorf("[HandleServerMsgPacket] player not exist")
}

func (w *World) Run() {
	startHTTPServer(w.httpPort, w.httpHandler, w.Config.HTTP.TLSCertFile, w.Config.HTTP.TLSKeyFile)
	go w.Server.Run()
	go w.playerManager.Run()
}

func (w *World) ForwardCrossZoneChatMsg(chatMsg *pbChat.SCCrossSrvChatMsg) {
	select {
	case w.crossZoneChatMsg <- chatMsg:
	default:
		logger.Error("[forwardCrossZoneChatMsg]  ZoneID:%v", chatMsg.ZoneID)
	}
}

func (w *World) ForwardSysMsg(sysMsg *pbChat.SCSystemMessage) {
	select {
	case w.systemMsgChan <- sysMsg:
	default:
		logger.Error("[forwardSysMsg] system message channel block:%v", sysMsg.Content)
	}
}
