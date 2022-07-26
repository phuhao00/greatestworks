package player

import "greatestworks/network/protocol/gen/messageId"

func (p *Player) HandlerRegister() {
	p.handlers[messageId.MessageId_CSAddFriend] = p.AddFriend
	p.handlers[messageId.MessageId_CSDelFriend] = p.DelFriend
	p.handlers[messageId.MessageId_CSSendChatMsg] = p.ResolveChatMsg
}
