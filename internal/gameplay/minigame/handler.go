package minigame

import (
	"errors"
	"github.com/phuhao00/greatestworks-proto/messageId"
	"github.com/phuhao00/network"
	"greatestworks/aop/module_router"
	"sync"
)

type Handler struct {
	Id messageId.MessageId
	Fn func(player Player, packet *network.Message)
}

var (
	handlers     []*Handler
	onceInit     sync.Once
	MinMessageId messageId.MessageId
	MaxMessageId messageId.MessageId //handle 的消息范围
)

func IsBelongToHere(id messageId.MessageId) bool {
	return id > MinMessageId && id < MaxMessageId
}

func GetHandler(id messageId.MessageId) (*Handler, error) {
	for _, handler := range handlers {
		if handler.Id == id {
			return handler, nil
		}
	}
	return nil, errors.New("not exist")
}

func RegisterHandler() {
	module_router.RegisterModuleMessageHandler(0, 0, nil)
}

func init() {
	onceInit.Do(func() {
		handlers[0] = &Handler{
			0,
			createGame,
		}
		handlers[1] = &Handler{
			1,
			leaveGame,
		}
	})
}

func createGame(player Player, message *network.Message) {

}

func leaveGame(player Player, message *network.Message) {

}
