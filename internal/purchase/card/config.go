package card

import "sync"

type Config struct {
	Id            uint32 `json:"id"`
	Category      uint32 `json:"category"`
	DailyReward   uint32 `json:"dailyReward"`
	LastDays      int32  `json:"lastDays"`
	RenewInterval int32  `json:"renewInterval"`
	Desc          string `json:"desc"`
}

var (
	configs sync.Map
)

func getCardConf(id uint32) *Config {
	value, ok := configs.Load(id)
	if !ok {
		return nil
	}
	return value.(*Config)
}

func getCardCategory(id uint32) Category {
	conf := getCardConf(id)
	if conf == nil {
		return NotDefine
	}
	return Category(conf.Category)
}
