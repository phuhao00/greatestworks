package event

import "sync"

type Normal struct {
	enum2OnEvent sync.Map
}

func (n *Normal) RegisterListener(e Enum, cb OnEvent) {
	n.enum2OnEvent.LoadOrStore(e, cb)
}

func (n *Normal) Dispatch(e Enum, params ...interface{}) {
	n.enum2OnEvent.Range(func(key, value any) bool {
		cb := value.(OnEvent)
		cb(params)
		return true
	})
}
