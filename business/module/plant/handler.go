package plant

import (
	"errors"
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"github.com/phuhao00/network"
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

func init() {
	onceInit.Do(func() {
		handlers[0] = &Handler{
			0,
			water,
		}
		handlers[1] = &Handler{
			1,
			spreadManure,
		}
		handlers[2] = &Handler{
			2,
			seed,
		}
		handlers[3] = &Handler{
			3,
			harvest,
		}
	})
}

// Water 浇水
func water(player Player, packet *network.Message) {

}

// spreadManure 施肥
func spreadManure(player Player, packet *network.Message) {

}

// seed 播种
func seed(player Player, packet *network.Message) {

}

// harvest 收割
func harvest(player Player, packet *network.Message) {

}
