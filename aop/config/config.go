package config

import "sync"

type Config struct {
	Path       string `yaml:"path"`
	Activity   string `yaml:"activity"`
	BattlePass string `yaml:"battlePass"`
	Pet        string `yaml:"pet"`
	Npc        string `yaml:"npc"`
	Plant      string `yaml:"plant"`
	Shop       string `yaml:"shop.proto"`
	Task       string `yaml:"task"`
	Skill      string `yaml:"skill"`
	Vip        string `yaml:"vipevent"`
	Building   string `yaml:"building"`
	Condition  string `yaml:"condition"`
	Synthetise string `yaml:"synthetise"`
	MiniGame   string `yaml:"miniGame"`
	Email      string `yaml:"email"`
}

var (
	config         Config
	onceInitConfig sync.Once
)

func GetConfig() Config {
	onceInitConfig.Do(func() {
		//todo 加载
	})
	return config
}

func BuildRealPath(base string, module string) string {
	return ""
}
