package pool

import "sync"

type Move struct {
}

var MovePool = sync.Pool{
	New: func() interface{} {
		msg := &Move{}
		return msg
	},
}
