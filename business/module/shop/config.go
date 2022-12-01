package shop

import "sync"

type Config struct {
	Id          uint32   `json:"id"`
	Desc        string   `json:"desc"`
	ItemIds     []uint32 `json:"itemIds"`
	Category    uint32   `json:"category"`    //商店类型
	RefreshTime []string `json:"refreshTime"` //刷新时间
}

var (
	Configs  sync.Map //配置
	onceInit sync.Once
)

func init() {
	onceInit.Do(func() {
		LoadConfigs()
	})
}

func LoadConfigs() {
	
}
