package action

import "github.com/looplab/fsm"

type Idle struct {
	event fsm.EventDesc
	cb    fsm.Callback
}

func (i *Idle) GetDesc() string {
	//TODO implement me
	panic("implement me")
}

func (i *Idle) GetEvent() fsm.EventDesc {
	//TODO implement me
	panic("implement me")
}

func (i *Idle) GetCb() fsm.Callback {
	//TODO implement me
	panic("implement me")
}
