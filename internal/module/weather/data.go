package weather

import "sync"

type Data struct {
}

var (
	weather         *Weather
	onceInitWeather sync.Once
)

func init() {
	onceInitWeather.Do(func() {
		weather = NewWeather()
	})
}

func GetWeather() *Weather {
	return weather
}
