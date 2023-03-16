package friend

import (
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"google.golang.org/protobuf/proto"
	"greatestworks/internal/event"
)

type IPlayer interface {
	Start()
	Stop()
	SendMsg(ID messageId.MessageId, message proto.Message)
	OnEvent(event event.IEvent)
}
