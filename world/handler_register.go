package world

import "greatestworks/network/protocol/gen/messageId"

func (mm *MgrMgr) HandlerRegister() {
	mm.Handlers[messageId.MessageId_CSLogin] = mm.CreatePlayer
	mm.Handlers[messageId.MessageId_CSLogin] = mm.UserLogin
}
