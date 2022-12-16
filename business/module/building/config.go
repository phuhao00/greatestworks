package building

import "sync"

type Config struct {
	Id              uint32
	Name            string
	Desc            string
	Cost            uint32
	UnlockCondition string
}

var (
	configs        sync.Map
	onceInitConfig sync.Once
)

func init() {
	onceInitConfig.Do(func() {
		//todo init configs
	})
}
