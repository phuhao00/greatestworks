package internal

// "greatestworks/internal/note/event" // TODO: 实现事件系统

// Manager 管理器接口定义
type Manager interface {
	OnStart()
	AfterStart()
	OnStop()
	AfterStop()
}

// Metrics 指标接口
type Metrics interface {
	GetName() string
	Description() string
	SetName(str string)
}

// DBAction 数据库操作接口
type DBAction interface {
	Load()
	Save()
}

// ConfigManagerAction 配置管理操作接口
type ConfigManagerAction interface {
	Load()
	Get(id uint32) interface{}
}

// Module 模块接口
type Module interface {
	OnEvent(c Character, event interface{}) // TODO: 实现event系统
	SetEventCategoryActive(eventCategory int)
	RegisterHandler()
	Manager
}
