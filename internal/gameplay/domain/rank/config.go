package rank

import (
	"fmt"
	"time"
)

type Config struct {
	ID          uint32 `json:"id"`
	Desc        string `json:"desc"`
	Name        string `json:"name"`
	Category    uint32 `json:"category"`
	SortType    uint32 `json:"sortType"`
	RefreshTime uint32 `json:"refreshTime"`
	Reward      uint32 `json:"reward"`
}

func (c *Config) getRankName(rankId uint32) (rankName string) {
	now := time.Now()
	rankName = fmt.Sprintf("molerank:%v:%04d%02d%02d", c.Category, now.Year(), now.Month(), now.Day())
	return rankName
}
