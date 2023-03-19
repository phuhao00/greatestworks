package main

import (
	"sync"
)

type GateWayList struct {
	sync.RWMutex
	endpoints *sync.Map
	ZoneId    int
	Levels    sync.Map //Level0    []string
}

func NewGatewayList(zoneId int) *GateWayList {
	return &GateWayList{endpoints: &sync.Map{}, ZoneId: zoneId}
}

func (l *GateWayList) removeLevel(id string) {
	if len(id) == 0 {
		return
	}
	l.Levels.Range(func(key, value any) bool {
		levelArr := value.([]string)
		for i, s := range levelArr {
			if id == s {
				copy(levelArr[i:], levelArr[i+1:])
				levelArr = levelArr[:len(levelArr)-1]
				l.Levels.Store(key, levelArr)
				return false
			}
		}
		return true
	})
}

func (l *GateWayList) removeEndPoint(id string) {
	if len(id) == 0 {
		return
	}
	l.removeLevel(id)
	l.endpoints.Delete(id)
}

func (l *GateWayList) clearInvalid() {

}

func (l *GateWayList) updateLevel() {

}

func (l *GateWayList) update() {

}

func (l *GateWayList) GetRecommend() {

}

func (l *GateWayList) exist() {

}

func (l *GateWayList) UpdateLocalWeight() {

}
