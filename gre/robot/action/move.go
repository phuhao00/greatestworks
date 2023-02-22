package action

import "github.com/looplab/fsm"

type Move struct {
	event fsm.EventDesc
	cb    fsm.Callback
}

func (m *Move) GetDesc() string {
	//TODO implement me
	panic("implement me")
}

func (m *Move) GetEvent() fsm.EventDesc {
	//TODO implement me
	panic("implement me")
}

func (m *Move) GetCb() fsm.Callback {
	//TODO implement me
	panic("implement me")
}
