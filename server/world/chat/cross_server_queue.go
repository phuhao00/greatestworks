package chat

import (
	goNsq "github.com/nsqio/go-nsq"
	"github.com/phuhao00/greatestworks-proto/chat"
	"google.golang.org/protobuf/proto"
	"greatestworks/aop/logger"
	"greatestworks/aop/nsq"
	"greatestworks/server/world/server"

	"sync"
)

var (
	csChatOnce sync.Once
	csChatIns  *CrossSrvChatHandler
)

type CrossSrvChatHandler struct {
	srvChatConsumer *goNsq.Consumer
}

func crossSrvChatGetMe() *CrossSrvChatHandler {
	csChatOnce.Do(func() {
		csChatIns = &CrossSrvChatHandler{}
	})
	return csChatIns
}

func (srvChat *CrossSrvChatHandler) initNsqHandler(channel string) {
	srvChat.srvChatConsumer = nsq.NewConsumer(nsq.ChatNSQ, nsq.PublicChat, channel, srvChat)
}

func (srvChat *CrossSrvChatHandler) HandleMessage(msg *goNsq.Message) error {
	chatMsg := &chat.SCCrossSrvChatMsg{}
	err := proto.Unmarshal(msg.Body, chatMsg)
	if err != nil {
		logger.Error("[HandleMessage] 收到聊天消息 消息解析错误")
		return nil
	}
	if chatMsg.SendTime < server.Oasis.StartTM {
		return nil
	}
	server.Oasis.ForwardCrossZoneChatMsg(chatMsg)
	logger.Debug("[HandleMessage] 收到聊天消息 地址:%v  消息:%v", msg.NSQDAddress, chatMsg.Content)
	return nil
}

func (srvChat *CrossSrvChatHandler) publishCrossSrvChatMsg(chatMsg interface{}) error {
	msgData, err := proto.Marshal(chatMsg.(proto.Message))
	if err != nil {
		logger.Error("[publishCrossSrvChatMsg] 聊天消息错误:", chatMsg)
		return err
	}
	err = nsq.PublishAsync(nsq.ChatNSQ, nsq.PublicChat, msgData, nil)
	if err != nil {
		logger.Error("[publishCrossSrvChatMsg] PublishAsync err ", err)
		return err
	}
	return nil
}

func (srvChat *CrossSrvChatHandler) stop() {
	if srvChat.srvChatConsumer != nil {
		srvChat.srvChatConsumer.Stop()
	}
}
