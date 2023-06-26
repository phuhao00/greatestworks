package config

import (
	"github.com/BurntSushi/toml"
	"greatestworks/aop/consul"
)

var (
	initConsulConfig consul.Config
)

func GetInitConsul() *consul.Config {
	_, err := toml.DecodeFile("./init_consul.toml", &initConsulConfig)
	if err != nil {
		return nil
	}
	return &initConsulConfig
}
