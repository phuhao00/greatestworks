package excel_object

import "sync"

type Template struct {
}

func (t *Template) Check() {
	//TODO implement me
	panic("implement me")
}

func (t *Template) Patch() {
	//TODO implement me
	panic("implement me")
}

type TemplateManager struct {
	Data sync.Map
}

func (t *TemplateManager) Get(key any) Object {
	//TODO implement me
	panic("implement me")
}

func (t *TemplateManager) Load(path string) error {
	//TODO implement me
	panic("implement me")
}

func (t *TemplateManager) LoadAfter() error {
	//TODO implement me
	panic("implement me")
}
