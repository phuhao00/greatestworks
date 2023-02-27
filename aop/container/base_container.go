package container

import "sync"

type BaseContainer struct {
	data sync.Map
}

func (b *BaseContainer) Save(query, update interface{}) {
	//TODO implement me
	panic("implement me")
}

func (b *BaseContainer) Set(tag string, val interface{}) {
	//TODO implement me
	panic("implement me")
}

func (b *BaseContainer) Get(tag string) interface{} {
	//TODO implement me
	panic("implement me")
}

func (b *BaseContainer) Add(vals interface{}) {
	//TODO implement me
	panic("implement me")
}

func (b *BaseContainer) Del(vals interface{}) {
	//TODO implement me
	panic("implement me")
}

func (b *BaseContainer) GetItem(val interface{}) interface{} {
	//TODO implement me
	panic("implement me")
}

func (b *BaseContainer) SetItem(val interface{}, items interface{}) {
	//TODO implement me
	panic("implement me")
}
