package honour

import (
	"sync"
	"time"
)

type Data struct {
	honours sync.Map
}

func (d *Data) Add(element *Element) {
	d.honours.Store(element.Id, element)
}

func (d *Data) CheckExist() bool {
	return false
}

type Element struct {
	Id          uint32
	HadActive   bool
	Removed     bool
	ExpiredTime int64
}

func (e *Element) SetExpiredTime(delta int64) {
	e.ExpiredTime = time.Now().Unix() + delta
}

func (e *Element) SetRemoved() {
	e.Removed = true
}

func (e *Element) SetHadActive() {
	e.HadActive = true
}
