package main

import "sync"

type Zone struct {
	GateWayList sync.Map
	WorldList   sync.Map
}

var (
	zone         *Zone
	zoneInitOnce sync.Once
)

func GetZone() *Zone {
	zoneInitOnce.Do(func() {
		zone = &Zone{
			GateWayList: sync.Map{},
			WorldList:   sync.Map{},
		}
	})
	zone.discovery()
	return zone
}

func (z *Zone) discovery() {
	z.discoveryGateWay()
	z.discoveryWorld()
}

func (z *Zone) discoveryGateWay() {

}

func (z *Zone) discoveryWorld() {

}

func (z *Zone) recommendZone() {

}
