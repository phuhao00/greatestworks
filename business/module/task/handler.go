package task

import (
	"errors"
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"github.com/phuhao00/network"
	"sync"
)

type Handler struct {
	Id messageId.MessageId
	Fn func(s *Manager, player Player, packet *network.Message)
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

func init() {
	onceInit.Do(func() {
		HandlerFriendRegister()
	})
}

func HandlerFriendRegister() {
	handlers[0] = &Handler{
		0,
		AcceptTask,
	}
	handlers[1] = &Handler{
		0,
		Submit,
	}
}

func AcceptTask(m *Manager, player Player, packet *network.Message) {

}

func Submit(m *Manager, player Player, packet *network.Message) {

}

func getCompleteCondition() {

}

func getReward(m *Manager) {

}

func checkUnlock(m *Manager) bool {
	return false
}
