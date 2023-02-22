package action

import "github.com/looplab/fsm"

type IAction interface {
	GetEvent() fsm.EventDesc
	GetCb() fsm.Callback
	GetDesc() string
}
