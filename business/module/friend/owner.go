package friend

import (
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"google.golang.org/protobuf/proto"
)

type Owner interface {
	Start()
	Stop()
	SendMsg(ID messageId.MessageId, message proto.Message)
}
