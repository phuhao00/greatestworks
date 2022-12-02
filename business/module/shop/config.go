package shop

import "sync"

type Config struct {
	Id          uint32   `json:"id"`
	Desc        string   `json:"desc"`
	ItemIds     []uint32 `json:"itemIds"`
	Weights     []uint32 `json:"weights"`
	Category    uint32   `json:"category"`    //商店类型
	RefreshTime []string `json:"refreshTime"` //刷新时间
}

var (
	Configs        sync.Map //配置
	onceInitConfig sync.Once
)

func init() {
	onceInitConfig.Do(func() {
		LoadConfigs()
	})
}

func LoadConfigs() {

}
