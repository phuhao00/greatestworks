package nsq

import (
	"github.com/nsqio/go-nsq"
	PbNsq "github.com/phuhao00/greatestworks-proto/nsq"
	"google.golang.org/protobuf/proto"
	"greatestworks/aop/logger"
	"greatestworks/server/world/server"
	"runtime/debug"
)

type MsgHandler func(cmd PbNsq.NsqCommand, data []byte)

type ComplexHandler struct {
	handlers map[PbNsq.NsqCommand]MsgHandler
}

func (cpx *ComplexHandler) HandleMessage(msg *nsq.Message) error {
	defer func() {
		if err := recover(); err != nil {
			stackMsg := string(debug.Stack())
			SendWarningMessage(server.Oasis.Config.HttpAddress, "online-nsq", err, stackMsg)
			logger.Error("[HandleMessage] player goroutine出错 err:%v stack:\n%v", err, stackMsg)
			return
		}
	}()
	request := PbNsq.ComplexMessage{}
	err := proto.Unmarshal(msg.Body, &request)
	if err != nil {
		logger.Error("[HandleMessage] 消息错误:", err)
		return nil
	}
	if request.Time < server.Oasis.StartTM {
		return nil
	}
	handler, ok := cpx.handlers[request.Cmd]
	if ok {
		handler(request.Cmd, request.Data)
		logger.Debug("[HandleMessage] 已处理NSQ消息 消息ID:%v", request.Cmd)
	} else {
		logger.Error("[HandleMessage] 未知的消息 cmd:%v", request.Cmd)
		return nil
	}
	return nil
}

func (cpx *ComplexHandler) registerHandler() {
	cpx.handlers[PbNsq.NsqCommand_MailSend] = cpx.emailHandler
}
