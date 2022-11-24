package pet

import (
	"errors"
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"github.com/phuhao00/network"
	"sync"
)

type Handler struct {
	Id messageId.MessageId
	Fn func(p Player, packet *network.Message)
}

var (
	handlers     []*Handler
	onceInit     sync.Once
	MinMessageId messageId.MessageId
	MaxMessageId messageId.MessageId //handle 的消息范围
)

func GetHandler(id messageId.MessageId) (*Handler, error) {

	if id > MinMessageId && id < MaxMessageId {
		return nil, errors.New("not in")
	}
	for _, handler := range handlers {
		if handler.Id == id {
			return handler, nil
		}
	}
	return nil, errors.New("not exist")
}

func init() {
	onceInit.Do(func() {
		HandlerPetRegister()
	})
}

func HandlerPetRegister() {
	handlers[0] = &Handler{
		0,
		AddPet,
	}
	handlers[1] = &Handler{
		0,
		DelPet,
	}
	handlers[2] = &Handler{
		0,
		UpdatePet,
	}
	handlers[3] = &Handler{
		0,
		GetPetInfo,
	}
	handlers[4] = &Handler{
		0,
		UpdatePictorial,
	}
	handlers[5] = &Handler{
		0,
		GetPictorialInfo,
	}
}

func AddPet(player Player, message *network.Message) {

}

func DelPet(player Player, message *network.Message) {

}

func UpdatePet(player Player, message *network.Message) {
	//todo  升级，升星（进化），合成，成长等级，设置状态

}

func GetPetInfo(player Player, message *network.Message) {
	//todo  单个获取，批量获取
}

// UpdatePictorial 更新图鉴
func UpdatePictorial(player Player, message *network.Message) {

}

// GetPictorialInfo 获取图鉴信息
func GetPictorialInfo(player Player, message *network.Message) {

}
