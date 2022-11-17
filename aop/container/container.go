package container

type IDelegate interface {
	Save(query, update interface{})
	Set(tag string, val interface{})
	Get(tag string) interface{}
}

type IContainer interface {
	IDelegate
	Add(vals interface{})
	Del(vals interface{})
	GetItem(val interface{}) interface{}
	SetItem(val interface{}, items interface{})
}
