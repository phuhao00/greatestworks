package action

import "github.com/looplab/fsm"

type Attack struct {
	event fsm.EventDesc
	cb    fsm.Callback
}

func (a *Attack) GetEvent() fsm.EventDesc {
	//TODO implement me
	panic("implement me")
}

func (a *Attack) GetCb() fsm.Callback {
	//TODO implement me
	panic("implement me")
}

func (a *Attack) GetDesc() string {
	panic("implement me")
}
