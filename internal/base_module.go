package internal

import (
	// "greatestworks/internal/infrastructure/module_router" // TODO: 实现模块路由
	// "greatestworks/internal/infrastructure/net/call" // TODO: 实现网络调用
	// "greatestworks/internal/note/event" // TODO: 实现事件系统

	"go.opentelemetry.io/otel/trace"
)

type BaseModule struct {
	ModuleName          string
	activeEventCategory map[int]bool
	tracer              trace.Tracer
	// methods             []call.MethodKey // TODO: 实现call包
}

func (b *BaseModule) OnEvent(c Character, event interface{}) {
	// TODO: 实现event处理
}

func (b *BaseModule) SetEventCategoryActive(eventCategory int) {
	b.activeEventCategory[eventCategory] = true
}

func NewBaseModule() *BaseModule {
	return &BaseModule{}
}

func (b *BaseModule) Get(id uint32) interface{} {
	return nil
}

func (b *BaseModule) Load() {

}

func (b *BaseModule) Save() {

}

func (b *BaseModule) GetName() string {
	return b.ModuleName
}

func (b *BaseModule) Description() string {
	return ""
}

func (b *BaseModule) SetName(str string) {
	b.ModuleName = str
}

func (b *BaseModule) OnStart() {

}

func (b *BaseModule) AfterStart() {

}

func (b *BaseModule) OnStop() {

}

func (b *BaseModule) AfterStop() {

}

func (b *BaseModule) RegisterHandler() {
	// TODO: 实现module_router
	// module_router.RegisterModuleMessageHandler(0, 0, nil)
}
