package weather

type Category uint16

const (
	Sunny Category = iota + 1
	Windy
	Rainy
	Snowy
	Cloudy
)

type Config struct {
	Year  uint32
	Month uint16
	Tags  [][]uint32 //day-{{time-time-weather},{time-time-weather}}
}
