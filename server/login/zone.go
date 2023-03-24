package main

import (
	"github.com/hashicorp/consul/api"
	"greatestworks/aop/consul"
	"greatestworks/server/login/config"
	"math/rand"
	"sync"
	"time"
)

type ZoneManager struct {
	GateWays sync.Map
	Worlds   sync.Map
}

var (
	zoneManager  *ZoneManager
	zoneInitOnce sync.Once
)

func GetZoneManager() *ZoneManager {
	zoneInitOnce.Do(func() {
		zoneManager = &ZoneManager{
			GateWays: sync.Map{},
			Worlds:   sync.Map{},
		}
	})
	zoneManager.discovery()
	return zoneManager
}

// getGateWay ..
func (z *ZoneManager) getGateWay(zoneId int) *GateWay {
	value, ok := z.GateWays.Load(zoneId)
	if ok && value != nil {
		return value.(*GateWay)
	}
	return nil
}

// getWorld ..
func (z *ZoneManager) getWorld(zoneId int) *World {
	value, ok := z.Worlds.Load(zoneId)
	if ok && value != nil {
		return value.(*World)
	}
	return nil
}

// discovery ..
func (z *ZoneManager) discovery() {
	z.discoveryGateWay()
	z.discoveryWorld()
}

// discoveryGateWay discovery gateway
func (z *ZoneManager) discoveryGateWay() {
	services, _, err := consul.QueryServices(config.GateWayServiceName)
	if err != nil {
		return
	}
	z.clearInvalidGateWay(services)

}

// discoveryWorld discover world
func (z *ZoneManager) discoveryWorld() {
	services, _, err := consul.QueryServices(config.WorldServiceName)
	if err != nil {
		return
	}
	z.clearInvalidWorld(services)

	for _, svcEntry := range services {
		if len(svcEntry.Service.Tags) < 1 {
			continue
		}
		performance := consul.GetPerformanceData(svcEntry.Service.Tags[0])

		maxUserCnt := GetServer().Conf.Me.MaxWorldPlayerNum
		// 最大人数
		if performance.MaxPlayerNum != 0 {
			maxUserCnt = uint32(performance.MaxPlayerNum)
		}
		endPoint := &WorldEndpoint{
			ZoneId:        performance.Zid,
			ID:            svcEntry.Service.ID,
			IP:            svcEntry.Service.Address,
			Port:          svcEntry.Service.Port,
			Name:          svcEntry.Service.Service,
			PlayerNum:     performance.PlayerNum,
			PIdx:          uint32(performance.PIdx),
			Max:           maxUserCnt,
			fakeServerCnt: int(GetServer().Conf.Me.PlayersServerCnt),
			BeginTM:       performance.BeginTM,
		}
		z.updateWorldEndPoint(endPoint)
	}
}

// clearInvalidGateWay clear invalid gateway
func (z *ZoneManager) clearInvalidGateWay(services []*api.ServiceEntry) {
	z.GateWays.Range(func(key, value any) bool {
		gateWay := value.(*GateWay)
		zoneId, isEmpty := gateWay.clearInvalid(services)
		if isEmpty {
			z.Worlds.Delete(zoneId)
		}
		return true
	})
}

// clearInvalidWorld  clear invalid world
func (z *ZoneManager) clearInvalidWorld(entries []*api.ServiceEntry) {
	z.Worlds.Range(func(key, value any) bool {
		world := value.(*World)
		zoneId, isEmpty := world.clearInvalid(entries)
		if isEmpty {
			z.Worlds.Delete(zoneId)
		}
		return true
	})
}

// updateWorldEndPoint  update world endPoint
func (z *ZoneManager) updateWorldEndPoint(endPoint *WorldEndpoint) {
	w := z.getWorld(endPoint.ZoneId)
	if w == nil {
		w = NewWorld(endPoint.ZoneId)
		w.updateOnline(endPoint)
		z.Worlds.Store(endPoint.ZoneId, w)

	} else {
		if w.ZoneId == endPoint.ZoneId {
			w.updateOnline(endPoint)
		}
	}
}

// existGateway check exist gateway endPoint
func (z *ZoneManager) existGateway(zoneId int, gatewayId string) (*config.EndPoint, bool) {
	gList := z.getGateWay(zoneId)
	if gList != nil {
		if gList.ZoneId == zoneId {
			return gList.GetEndPoint(gatewayId), true
		}
	}
	return nil, false
}

// worldMetrics world metrics
func (z *ZoneManager) worldMetrics(endPoint *WorldEndpoint) int32 {
	if endPoint == nil {
		return config.CloseStatus
	} else if endPoint.Max == 0 {
		return config.OKStatus
	} else {
		playerNum := endPoint.PlayerNum
		if int(endPoint.PIdx) < endPoint.fakeServerCnt && endPoint.initFakePlayerNum > endPoint.PlayerNum {
			now := time.Now().Unix()
			if int32(now-endPoint.BeginTM) < GetServer().Conf.Me.PlayerNumHour*config.HoursSeconds {
				playerNum = endPoint.initFakePlayerNum
			} else if int32(now-endPoint.BeginTM) > GetServer().Conf.Me.PlayerNumHour*config.HoursSeconds {
				endPoint.initFakePlayerNum = 0
			}
		}
		var stat int32
		rate := float32(playerNum) / float32(endPoint.Max)
		if rate < config.EmptyRatio {
			stat = config.EmptyStatus
		} else if rate < config.BusyRatio {
			stat = config.OKStatus
		} else {
			stat = config.FullStatus
		}
		return stat
	}
}

// recommendZone recommend zoneManager
func (z *ZoneManager) recommendZone() int {
	recZoneId := 0
	z.Worlds.Range(func(zoneId, value interface{}) bool {
		emptyCnt := 0
		okCnt := 0
		w := value.(*World)
		w.endPoints.Range(func(sid, value interface{}) bool {
			endpoint := value.(*WorldEndpoint)
			stat := z.worldMetrics(endpoint)
			if stat == config.EmptyStatus {
				emptyCnt++
			} else if stat == config.OKStatus {
				okCnt++
			}
			return true
		})

		if emptyCnt > config.RecommendWorldMaxCnt {
			recZoneId = zoneId.(int)
			return false
		}

		if recZoneId == 0 && okCnt > 0 {
			recZoneId = zoneId.(int)
			return false
		}
		return true
	})

	if recZoneId == 0 {
		var zoneIds []int
		z.Worlds.Range(func(zoneId, value interface{}) bool {
			zoneIds = append(zoneIds, zoneId.(int))
			return true
		})
		if len(zoneIds) > 0 {
			rnd := rand.Int()
			if rnd < 0 {
				rnd = -rnd
			}
			recZoneId = zoneIds[rnd%len(zoneIds)]
		}
	}
	return recZoneId
}
