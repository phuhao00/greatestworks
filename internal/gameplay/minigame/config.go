package minigame

type Config struct {
	Id         uint32 `json:"Id"`
	Category   uint16 `json:"Category"`
	OpenTime   int64  `json:"OpenTime"`
	CloseTime  int64  `json:"CloseTime"`
	Reward     uint32 `json:"Reward"`
	Limit      string `json:"Limit"`
	MaxPlayer  uint8  `json:"MaxPlayer"`
	MiniPlayer uint8  `json:"MiniPlayer"`
}
