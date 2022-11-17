package family

import (
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"google.golang.org/protobuf/proto"
)

type Owner interface {
	Start()
	Stop()
	BroadcastMsg(ids []uint64, msgId messageId.MessageId, msg proto.Message)
}
