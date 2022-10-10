package template

import "time"

type HappyNewYear struct {
	Id        uint32
	StartTime time.Time
	EndTime   time.Time
}

func (y *HappyNewYear) Init(conf Conf) *HappyNewYear {
	return &HappyNewYear{}
}
