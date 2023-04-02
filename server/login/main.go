package main

import (
	"greatestworks/aop/consul"
	"greatestworks/aop/logger"
	"greatestworks/server/login/config"
)

func main() {
	//
	err := consul.InitConsul(nil)
	if err != nil {
		return
	}
	cfg := &config.Config{}
	consul.LoadJSONFromConsulKV(consul.GetConsulConfigName(), cfg)
	logger.SetLogging(&logger.LoggingSetting{})
	//init mongo
	config.QueryToGateWayRatio = cfg.Me.QueryGateWayRatio
	//todo name mod init
	//todo nsq init
	//todo redis init
	//todo dirty filter init
	//todo token load

	server := GetServer()
	server.Loop()
}
