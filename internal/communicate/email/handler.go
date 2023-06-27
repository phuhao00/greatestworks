package email

import (
	"errors"
	"github.com/phuhao00/greatestworks-proto/messageId"
	"github.com/phuhao00/network"
	"sync"
)

type Handler struct {
	Id messageId.MessageId
	Fn func(player IPlayer, packet *network.Message)
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
		handlers[0] = &Handler{
			0,
			readMail,
		}
		handlers[1] = &Handler{
			1,
			deleteMail,
		}
	})
}

func readMail(player IPlayer, message *network.Message) {

}

func deleteMail(player IPlayer, message *network.Message) {

}
