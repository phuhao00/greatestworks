package weather

import (
	"sync"
	"time"
)

type Weather struct {
	CurStatus Category //当前天气状态
	configs   sync.Map
	openTime  time.Time
}

func NewWeather() *Weather {
	return &Weather{
		CurStatus: 0,
		configs:   sync.Map{},
		openTime:  time.Time{},
	}
}

func (w *Weather) Update() {
	ti := time.NewTicker(time.Second)
	for {
		select {
		case <-ti.C:
			w.calcWeather()
		}
	}
}

func (w *Weather) GetWeatherStatus() Category {
	return w.CurStatus
}

func (w *Weather) calcWeather() {
	//todo load from configs
	if true {
		//todo change weather status
	}
}
