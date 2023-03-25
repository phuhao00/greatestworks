package main

import "github.com/phuhao00/greatestworks-proto/messageId"

func (w *World) HandlerRegister() {
	w.Handlers[messageId.MessageId_CSCreatePlayer] = w.CreatePlayer
	w.Handlers[messageId.MessageId_CSLogin] = w.UserLogin
}
