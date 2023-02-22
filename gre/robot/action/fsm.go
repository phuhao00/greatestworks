package action

import (
	"github.com/looplab/fsm"
)

func NewFsmWrap() *FsmWrap {
	f := &FsmWrap{fsm.FSM{}, make([]IAction, 0), ""}
	return f
}

type FsmWrap struct {
	fsm.FSM
	actions    []IAction
	InitAction string
}

func (f *FsmWrap) AddAction(action IAction) {
	f.actions = append(f.actions, action)
}

func (f *FsmWrap) Init() {
	events, cbs := func() (evs []fsm.EventDesc, cbs map[string]fsm.Callback) {
		evs = make([]fsm.EventDesc, 0, len(f.actions))
		cbs = make(map[string]fsm.Callback, len(f.actions))
		for _, action := range f.actions {
			evs = append(evs, action.GetEvent())
			cbs[action.GetDesc()] = action.GetCb()
		}
		return evs, cbs
	}()
	fsm.NewFSM(f.InitAction, events, cbs)
}
