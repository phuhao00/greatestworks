package task

type Task interface {
	Accept(config *Config)
	Finish()
	TargetDoneCallBack()
}

type Base struct {
}

func (b *Base) Accept(config *Config) {

}

func (b *Base) Finish() {

}

func (b *Base) TargetDoneCallBack() {

}
