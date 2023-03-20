package main

import (
	"greatestworks/server/login/config"
	"sync"
)

type GateWayList struct {
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

func (l *GateWayList) addLevel(id string, lv int) {
	if len(id) == 0 {
		return
	}
	l.Levels.Range(func(key, value any) bool {
		if lv != key.(int) {
			return true
		}
		levelArr := value.([]string)
		for _, s := range levelArr {
			if id == s {
				levelArr = append(levelArr, id)
				l.Levels.Store(key, levelArr)
				return false
			}
		}
		return true
	})
}

func (l *GateWayList) name() {

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

func (l *GateWayList) update(ep config.EndPoint) {
	if l == nil {
		return
	}
	if ep.Weights < config.LEVEL0 {
		l.removeLevel(ep.ID)
		l.addLevel(ep.ID, config.LEVEL0)
	} else if ep.Weights < config.LEVEL1 {
		l.removeLevel(ep.ID)
		l.addLevel(ep.ID, config.LEVEL1)
	} else if ep.Weights < config.LEVEL2 {
		l.removeLevel(ep.ID)
		l.addLevel(ep.ID, config.LEVEL2)
	} else if ep.Weights < config.LEVEL3 {
		l.removeLevel(ep.ID)
		l.addLevel(ep.ID, config.LEVEL3)
	}
}

func (l *GateWayList) GetRecommend() {

}

func (l *GateWayList) exist() {

}

func (l *GateWayList) UpdateLocalWeight() {

}
