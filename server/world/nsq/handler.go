package nsq

import (
	"github.com/phuhao00/greatestworks-proto/mail"
	PbNsq "github.com/phuhao00/greatestworks-proto/nsq"
	"google.golang.org/protobuf/proto"
	"greatestworks/aop/logger"
)

func (cpx *ComplexHandler) emailHandler(cmd PbNsq.NsqCommand, data []byte) {
	m := mail.MailInfo{}
	err := proto.Unmarshal(data, &m)
	if err != nil {
		logger.Error("[ComplexHandler] 消息错误 err:", err)
		return
	}
	logger.Debug("[ComplexHandler]  %v", m)
}
