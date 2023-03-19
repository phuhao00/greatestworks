package config

type GateWay struct {
	ZoneId  int
	ID      string
	IP      string
	Port    int
	Name    string
	Weights int32
	InnerIP string
}
