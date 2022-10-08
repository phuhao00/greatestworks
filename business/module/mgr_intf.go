package module

// MgrInterface 管理器接口定义
type MgrInterface interface {
	OnStart()
	AfterStart()
	OnStop()
	AfterStop()
}

type Metrics interface {
	GetName() string
	SetName(str string)
}
