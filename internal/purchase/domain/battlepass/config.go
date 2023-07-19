package battlepass

import "time"

type Config struct {
	Id        uint32
	Score     []uint64
	Free      []uint32
	Gold      []uint32
	Price     uint64
	Desc      string
	OpenTime  time.Time
	CloseTime time.Time
}
type CardCategory uint8

const (
	Free CardCategory = iota + 1
	Gold
)
