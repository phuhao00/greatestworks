package chat

import (
	"fmt"
	goNsq "github.com/nsqio/go-nsq"
	"github.com/phuhao00/greatestworks-proto/chat"
	"google.golang.org/protobuf/proto"
	"greatestworks/aop/fn"
	"greatestworks/aop/logger"
	"greatestworks/aop/nsq"
	"greatestworks/server/world/server"
	"strings"
	"sync"
)

const (
	NoneMsg   uint32 = 0 //
	SysMsgEnd uint32 = 5 //
)

var (
	sysMsgOnce sync.Once
	sysMsgIns  *SystemMsgHandler
)

type SystemMsgHandler struct {
	sysMsgConsumer *goNsq.Consumer
}

func systemMsgGetMe() *SystemMsgHandler {
	sysMsgOnce.Do(func() {
		sysMsgIns = &SystemMsgHandler{}
	})
	return sysMsgIns

}

func (sysMsgHandler *SystemMsgHandler) initNsqHandler(channel string) {
	err := nsq.CreateTopic(nsq.LogicNSQ, nsq.SystemMsg)
	if err != nil {
		logger.Error("[initNsqHandler]err:%v", err)
	}
	sysMsgHandler.sysMsgConsumer = nsq.NewConsumer(nsq.LogicNSQ, nsq.SystemMsg, channel, sysMsgHandler)
}

func (sysMsgHandler *SystemMsgHandler) HandleMessage(msg *goNsq.Message) error {
	sysMsg := &chat.SCSystemMessage{}
	err := proto.Unmarshal(msg.Body, sysMsg)
	if err != nil {
		logger.Error("err:%v", err)
		return nil
	}

	if sysMsg.SendTime < server.Oasis.StartTM {
		logger.Debug("[HandleMessage] 收到系统消息, 老于本服启动时间，被忽略, Modules:%v, Channels:%v, Content:%v.",
			sysMsg.Modules, sysMsg.Channels, sysMsg.Content)
		return nil
	} else {
		logger.Debug("[HandleMessage] 收到系统消息, Modules:%v, Channels:%v, Content:%v",
			sysMsg.Modules, sysMsg.Channels, sysMsg.Content)
	}

	if len(sysMsg.Modules) > 0 {
		if strings.ToLower(sysMsg.Modules) == "all" {
			server.Oasis.ForwardSysMsg(sysMsg)
		} else {
			pIds := fn.SplitStringToInt32Slice(sysMsg.Modules, ",")
			for _, pid := range pIds {
				if pid >= 0 && pid == int32(server.Oasis.Pid) {
					server.Oasis.ForwardSysMsg(sysMsg)
					break
				}
			}
		}
	} else {
		server.Oasis.ForwardSysMsg(sysMsg)
	}

	return nil
}

func (sysMsgHandler *SystemMsgHandler) publishSysMsg(sysMsg interface{}) error {
	msgData, err := proto.Marshal(sysMsg.(proto.Message))
	if err != nil {
		logger.Error("[publishSysMsg Marshal err:", err.Error())
		return err
	}
	err = nsq.PublishAsync(nsq.LogicNSQ, nsq.SystemMsg, msgData, nil)
	if err != nil {
		logger.Error("[publishSysMsg] PublishAsync err ", err)
		return err
	}
	return nil
}

func (sysMsgHandler *SystemMsgHandler) sendSysMsg(msgType uint32, format string, a ...interface{}) error {
	if msgType <= NoneMsg || msgType >= SysMsgEnd {
		err := fmt.Errorf("system message Type:%v", msgType)
		return err
	}
	content := fmt.Sprintf(format, a...)
	sysMsg := &chat.SCSystemMessage{MsgType: msgType, Content: content}
	return sysMsgHandler.publishSysMsg(sysMsg)
}

func (sysMsgHandler *SystemMsgHandler) stop() {
	if sysMsgHandler.sysMsgConsumer != nil {
		sysMsgHandler.sysMsgConsumer.Stop()
	}
}
