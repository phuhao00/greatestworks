package action

import "github.com/looplab/fsm"

type Attack struct {
	event fsm.EventDesc
	cb    fsm.Callback
}
