package player

import (
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"github.com/phuhao00/network"
	"greatestworks/business/module/friend"
	"sync"
)

type Handler func(p *Player, packet *network.Message)

var (
	handlers map[messageId.MessageId]Handler
	onceInit sync.Once
)

func init() {
	onceInit.Do(func() {
		HandlerChatRegister()
		friend.HandlerFriendRegister()
	})
}
