package module_router

import (
	"github.com/phuhao00/greatestworks-proto/messageId"
	"github.com/phuhao00/greatestworks-proto/module"
)

type ModuleMessageHandler func(messageId uint64, data []byte)

var (
	Module2MessageId2Handler map[module.Module]map[messageId.MessageId]ModuleMessageHandler
)

func RegisterModuleMessageHandler(moduleId module.Module, msgId messageId.MessageId, handler ModuleMessageHandler) {
	if Module2MessageId2Handler == nil {
		Module2MessageId2Handler = make(map[module.Module]map[messageId.MessageId]ModuleMessageHandler)
	}
	if Module2MessageId2Handler[moduleId] == nil {
		Module2MessageId2Handler[moduleId] = make(map[messageId.MessageId]ModuleMessageHandler)
	}
	if Module2MessageId2Handler[moduleId][msgId] != nil {
		panic("[RegisterModuleMessageHandler] repeated register")
	}
	Module2MessageId2Handler[moduleId][msgId] = handler
}

func GetModuleHandler(moduleId module.Module, messageId messageId.MessageId) ModuleMessageHandler {
	message2Handler, ok := Module2MessageId2Handler[moduleId]
	if !ok {
		return nil
	}
	handler, exist := message2Handler[messageId]
	if !exist {
		return nil
	}
	return handler
}
