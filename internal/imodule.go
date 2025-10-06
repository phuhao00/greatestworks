package internal

// "greatestworks/internal/note/event" // TODO: 实现事件系统

// IManager 管理器接口定义
type IManager interface {
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

// DBAction 加载、 存储DB
type DBAction interface {
	Load()
	Save()
}

type ConfigMgrAction interface {
	Load()
	Get(id uint32) interface{}
}

type IModule interface {
	OnEvent(c Character, event interface{}) // TODO: 实现event系统
	SetEventCategoryActive(eventCategory int)
	RegisterHandler()
	IManager
}
