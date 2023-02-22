package action

import "github.com/looplab/fsm"

type Chase struct {
	event fsm.EventDesc
	cb    fsm.Callback
}

func (c *Chase) GetDesc() string {
	//TODO implement me
	panic("implement me")
}

func (c *Chase) GetEvent() fsm.EventDesc {
	//TODO implement me
	panic("implement me")
}

func (c *Chase) GetCb() fsm.Callback {
	//TODO implement me
	panic("implement me")
}
