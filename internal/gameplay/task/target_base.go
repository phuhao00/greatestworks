package task

type BaseTarget struct {
}

func NewBaseTarget() *BaseTarget {
	return &BaseTarget{}
}

func (b *BaseTarget) GetProgress() {
	//TODO implement me
	panic("implement me")
}

func (b *BaseTarget) GetTotalProgress() {
	//TODO implement me
	panic("implement me")
}

func (b *BaseTarget) CheckDone() {
	//TODO implement me
	panic("implement me")
}

func (b *BaseTarget) OnAccept() {
	//TODO implement me
	panic("implement me")
}

func (b *BaseTarget) OnEvent(i interface{}) {
	//TODO implement me
	panic("implement me")
}

func (b *BaseTarget) ConfigVerify() {
	//TODO implement me
	panic("implement me")
}

func (b *BaseTarget) Notify() {

}
