package world

import "google.golang.org/protobuf/proto"

func (w *World) BroadcastSystemMsg(message proto.Message) {
	//TODO implement me
	panic("implement me")
}

func (w *World) BroadcastOnlineChatMsg(message proto.Message) {
	//TODO implement me
	panic("implement me")
}

func (w *World) BroadcastCrossZoneChatMsg(message proto.Message) {
	//TODO implement me
	panic("implement me")
}

func (w *World) BroadcastZoneChatMsg(message proto.Message) {
	//TODO implement me
	panic("implement me")
}

func (w *World) BroadcastCrossSrvChatMsg(message proto.Message) {
	//TODO implement me
	panic("implement me")
}

func (w *World) SyncOfflineOnlineChatMsg() []proto.Message {
	//TODO implement me
	panic("implement me")
}
