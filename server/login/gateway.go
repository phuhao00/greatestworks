package main

import (
	"github.com/hashicorp/consul/api"
	"greatestworks/aop/logger"
	"greatestworks/server/login/config"
	"math/rand"
	"sync"
)

type GateWay struct {
	endpoints *sync.Map
	ZoneId    int
	Levels    sync.Map //Level0    []string
}

func NewGatewayList(zoneId int) *GateWay {
	return &GateWay{endpoints: &sync.Map{}, ZoneId: zoneId}
}

func (l *GateWay) removeLevel(id string) {
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

func (l *GateWay) addLevel(id string, lv int) {
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

func (l *GateWay) removeEndPoint(id string) {
	if len(id) == 0 {
		return
	}
	l.removeLevel(id)
	l.endpoints.Delete(id)
}

func (l *GateWay) clearInvalid(services []*api.ServiceEntry) (zoneId int, isEmpty bool) {
	totalCnt := 0
	var del = make([]string, 0)
	l.endpoints.Range(func(key, value any) bool {
		e := value.(*config.EndPoint)
		if e.ZoneId != l.ZoneId {
			return true
		}
		totalCnt++
		exist := false
		for _, service := range services {
			if len(e.ID) > 0 && e.ID == service.Service.ID {
				exist = true
			}
		}
		if !exist {
			del = append(del, e.ID)
		}
		return true
	})
	delCnt := 0
	for _, s := range del {
		l.removeEndPoint(s)
		delCnt++
	}
	return l.ZoneId, delCnt == totalCnt
}

func (l *GateWay) update(ep *config.EndPoint) {
	if l == nil || ep.ZoneId != l.ZoneId {
		logger.Error("GatewayList's zoneManager id error.")
		return
	}
	needUp := false
	value, ok := l.endpoints.Load(ep.ID)
	if value != nil && ok {
		if value.(*config.EndPoint).Weights != ep.Weights {
			needUp = true
			value.(*config.EndPoint).Weights = ep.Weights
		}
	} else {
		l.endpoints.Store(ep.ID, ep)
		needUp = true
	}
	if !needUp {
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

func (l *GateWay) GetRecommend() *config.EndPoint {
	r := rand.Int()
	var ret *config.EndPoint
	l.Levels.Range(func(key, value any) bool {
		levelArr := value.([]string)
		addr := levelArr[r%len(levelArr)]
		endpoint, ok := l.endpoints.Load(addr)
		if ok && endpoint != nil {
			ret = endpoint.(*config.EndPoint)
			return false
		}
		return true
	})
	if ret != nil {
		return ret
	}
	return &config.EndPoint{IP: "127.0.0.1", Port: 10001}
}

func (l *GateWay) GetEndPoint(addr string) *config.EndPoint {
	var ret *config.EndPoint
	value, ok := l.endpoints.Load(addr)
	if ok && value != nil {
		l.Levels.Range(func(key, value any) bool {
			arr := value.([]string)
			for _, s := range arr {
				if addr == s {
					ret = value.(*config.EndPoint)
					return false
				}
			}
			return true
		})
	}
	return ret
}

func (l *GateWay) UpdateLocalWeight(endPoint *config.EndPoint) {
	if endPoint != nil {
		value, ok := l.endpoints.Load(endPoint.ID)
		if value != nil && ok {
			if config.QueryToGateWayRatio < 1 {
				config.QueryToGateWayRatio = 1
			}
			value.(*config.EndPoint).Weights += config.QueryToGateWayRatio
			l.update(value.(*config.EndPoint))
		}
	}
}

func (l *GateWay) IsExist(addr string) (*config.EndPoint, bool) {
	var ret *config.EndPoint
	if endpoint, ok := l.endpoints.Load(addr); ok && endpoint != nil {
		l.Levels.Range(func(key, value any) bool {
			arr := value.([]string)
			for _, s := range arr {
				if addr == s {
					ret = value.(*config.EndPoint)
					return false
				}
			}
			return true
		})
	}
	return ret, false
}
