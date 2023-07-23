package card

import "sync"

type Config struct {
	Category      uint32 `json:"category"`
	DailyReward   uint32 `json:"dailyReward"`
	LastDays      int32  `json:"lastDays"`
	RenewInterval int32  `json:"renewInterval"`
	BuyReward     uint32 `json:"buyReward"`
	Desc          string `json:"desc"`
}

var (
	configs sync.Map
)

func getCardConf(category Category) *Config {
	value, ok := configs.Load(category)
	if !ok {
		return nil
	}
	return value.(*Config)
}
