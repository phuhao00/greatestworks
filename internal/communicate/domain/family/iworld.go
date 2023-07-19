package family

import (
	"github.com/phuhao00/greatestworks-proto/messageId"
	"google.golang.org/protobuf/proto"
)

type IWorld interface {
	Start()
	Stop()
	BroadcastMsg(ids []uint64, msgId messageId.MessageId, msg proto.Message)
}
