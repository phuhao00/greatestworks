package module_router

type ModuleMessageHandler func(messageId uint64, data []byte)

var (
	Module2MessageId2Handler map[uint16]map[uint64]ModuleMessageHandler
)

func RegisterModuleMessageHandler(moduleId uint16, messageId uint64, handler ModuleMessageHandler) {
	if Module2MessageId2Handler == nil {
		Module2MessageId2Handler = make(map[uint16]map[uint64]ModuleMessageHandler)
	}
	if Module2MessageId2Handler[moduleId] == nil {
		Module2MessageId2Handler[moduleId] = make(map[uint64]ModuleMessageHandler)
	}
	Module2MessageId2Handler[moduleId][messageId] = handler
}

func GetModuleHandler(moduleId uint16, messageId uint64) ModuleMessageHandler {
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
