package task

import (
	"github.com/phuhao00/greatestworks-proto/messageId"
	"google.golang.org/protobuf/proto"
)

type Player interface {
	SendMsg(ID messageId.MessageId, message proto.Message)
	GetTaskData() *Data
}
