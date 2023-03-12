package module

import "greatestworks/aop/event"

type BaseModule struct {
	ModuleName          string
	activeEventCategory map[int]bool
}

func (b *BaseModule) OnEvent(c Character, event event.IEvent) {

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
