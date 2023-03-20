package config

import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Config
// 使用指针，考虑到GM修改
type Config struct {
	Log        *Log //
	Mongo      *Mongo
	Redis      *Redis
	RabbitMq   *RabbitMq
	WhiteList  *WhiteList
	ThirdParty *ThirdParty
	Me         *Me
	Consul     *Consul
	Etcd       *Etcd
	GateWays   []*GateWay
}

func Deserialize(str string) *Config {
	ret := &Config{}
	json.UnmarshalFromString(str, ret)
	return ret
}
