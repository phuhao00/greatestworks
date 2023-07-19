package friend

import (
	"github.com/phuhao00/greatestworks-proto/messageId"
	"google.golang.org/protobuf/proto"
	"greatestworks/internal/note/event"
)

type IPlayer interface {
	Start()
	Stop()
	SendMsg(ID messageId.MessageId, message proto.Message)
	OnEvent(event event.IEvent)
}
