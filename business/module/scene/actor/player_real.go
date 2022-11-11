package actor

import "google.golang.org/protobuf/proto"

type PlayerReal interface {
	SendMsg(message proto.Message)
}
