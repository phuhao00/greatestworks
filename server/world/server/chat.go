package server

import (
	"google.golang.org/protobuf/proto"
)

func (w *World) BroadcastSystemMsg(message proto.Message) {
}

func (w *World) BroadcastOnlineChatMsg(message proto.Message) {
}

func (w *World) BroadcastCrossZoneChatMsg(message proto.Message) {
}

func (w *World) BroadcastZoneChatMsg(message proto.Message) {
}

func (w *World) BroadcastCrossSrvChatMsg(message proto.Message) {
}

func (w *World) SyncOfflineOnlineChatMsg() []proto.Message {
	return nil
}
