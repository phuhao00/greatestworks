package excel_object

type Drop struct {
	Id    uint32
	Items []uint32
}

func (d Drop) Check() {
	//TODO implement me
	panic("implement me")
}

func (d Drop) Patch() {
	//TODO implement me
	panic("implement me")
}

type DropManager struct {
}

func (d DropManager) Get(key any) Object {
	//TODO implement me
	panic("implement me")
}

func (d DropManager) Load(path string) error {
	//TODO implement me
	panic("implement me")
}

func (d DropManager) LoadAfter() error {
	//TODO implement me
	panic("implement me")
}
