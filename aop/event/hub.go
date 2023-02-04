package event

type OnEvent func(params ...interface{})

type Hub interface {
	RegisterListener(e Enum, cb OnEvent)
	Dispatch(e Enum, params ...interface{})
}
