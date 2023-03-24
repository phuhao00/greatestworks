package main

import (
	"github.com/hashicorp/consul/api"
	"greatestworks/server/login/config"
	"sync"
	"time"
)

type WorldEndpoint struct {
	ZoneId            int
	ID                string
	IP                string
	Port              int
	Name              string
	PlayerNum         int32
	PIdx              uint32
	Max               uint32
	BeginTM           int64
	initFakePlayerNum int32
	fakeServerCnt     int
}

type World struct {
	endPoints sync.Map
	ZoneId    int
}

var (
	world *World
)

func NewWorld(zoneId int) *World {
	return &World{
		endPoints: sync.Map{},
		ZoneId:    zoneId,
	}
}

func (w *World) removeEndpoint(id string) {
	if len(id) == 0 {
		return
	}
	w.endPoints.Delete(id)
}

func (w *World) clearInvalid(entries []*api.ServiceEntry) (zoneId int, isEmpty bool) {
	totalCnt := 0
	var del []string
	w.endPoints.Range(func(sid, value interface{}) bool {
		endpoint := value.(*WorldEndpoint)
		if endpoint.ZoneId != w.ZoneId {
			return true
		}
		totalCnt++

		exist := false
		for _, svcEntry := range entries {
			if len(endpoint.ID) > 0 && endpoint.ID == svcEntry.Service.ID {
				exist = true
				break
			}
		}
		if !exist {
			del = append(del, endpoint.ID)
		}
		return true
	})

	delCnt := 0
	for _, id := range del {
		w.removeEndpoint(id)
		delCnt++
	}
	return w.ZoneId, totalCnt == delCnt
}

func (w *World) updateOnline(enPoint *WorldEndpoint) {

	if enPoint.ZoneId != w.ZoneId {
		return
	}

	if ept, ok := w.endPoints.Load(enPoint.ID); ok {
		endpoint := ept.(*WorldEndpoint)

		if endpoint.PlayerNum != enPoint.PlayerNum {
			endpoint.PlayerNum = enPoint.PlayerNum
		}

		if endpoint.Max != enPoint.Max {
			endpoint.Max = enPoint.Max
		}

		if int(enPoint.PIdx) < enPoint.fakeServerCnt && enPoint.initFakePlayerNum <= 0 {
			now := time.Now().Unix()
			if int32(now-endpoint.BeginTM) < GetServer().Conf.Me.PlayerNumHour*config.HoursSeconds {
				enPoint.initFakePlayerNum = int32(float32(enPoint.Max)*config.WorldMaxCoefficient) - GetServer().Conf.Me.PlayersDeltaCnt
			}
		}
	} else {
		if int(enPoint.PIdx) < enPoint.fakeServerCnt {
			enPoint.initFakePlayerNum = int32(float32(enPoint.Max)*config.WorldMaxCoefficient) - GetServer().Conf.Me.PlayersDeltaCnt
		}
		w.endPoints.Store(enPoint.ID, enPoint)
	}
}
