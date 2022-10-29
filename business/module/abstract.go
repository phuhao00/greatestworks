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
	Description() string
	SetName(str string)
}

//DBAction 加载、 存储DB
type DBAction interface {
	Load()
	Save()
}

type ConfigMgrAction interface {
	Load()
	Get(id uint32) interface{}
}

type Abstract interface {
}
