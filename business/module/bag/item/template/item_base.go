package template

import "time"

type ItemBase struct {
	Id             uint32 `json:"id"`
	Num            int64  `json:"num"`
	LastChangeTime int64  `json:"lastChangeTime"`
	UseTime        int64  `json:"useTime"`
}

func (i *ItemBase) Add(delta int64) {
	i.Num += delta
	i.LastChangeTime = time.Now().Unix()
}

func (i *ItemBase) Delete(delta int64) {
	i.Num -= delta
	i.LastChangeTime = time.Now().Unix()
}

func (i *ItemBase) GetNum() int64 {
	return i.Num
}

func (i *ItemBase) GetId() uint32 {
	return i.Id
}
