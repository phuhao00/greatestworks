package world

import "github.com/phuhao00/greatestworks-proto/gen/messageId"

func (w *World) HandlerRegister() {
	w.Handlers[messageId.MessageId_CSCreatePlayer] = w.CreatePlayer
	w.Handlers[messageId.MessageId_CSLogin] = w.UserLogin
}
