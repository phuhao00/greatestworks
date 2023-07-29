package main

import (
	"context"
	"flag"
	"github.com/phuhao00/spoor"
	"greatestworks/aop/consul"
	"greatestworks/aop/fn"
	"greatestworks/aop/logger"
	"greatestworks/aop/redis"
	"greatestworks/server/gateway/config"
	"greatestworks/server/gateway/server"
	"strconv"
)

var (
	pid = flag.Int("pid", 1, "the same process number")
)

func main() {

	flag.Parse()

	err := consul.InitConsul(config.GetInitConsul())
	if err != nil {
		logger.Error("[main.go] Consul初始化失败error:%v", err)
		return
	}

	privateIP, err := fn.GetPrivateIPv4()

	if err != nil {
		logger.Error("[main.go] Get local ip error ", err)
		return
	}

	confName := "consul:" + privateIP + "-" + fn.GetUser() + "-" + "gateway.json"

	var cfg *config.Config

	consul.LoadJSONFromConsulKV(confName, &cfg)

	if cfg == nil {
		logger.Error("[main.go]load config fail!!!")
		return
	}
	logLevel, err := spoor.ParseLogLevel(cfg.Log.LogLevel)
	if err != nil {
		panic(err)
	}
	logSetting := &logger.LoggingSetting{
		Dir:    cfg.Log.LogPath + "_" + strconv.Itoa(*pid) + cfg.Log.LogFile,
		Level:  int(logLevel),
		Prefix: "[gateway]",
	}
	logger.SetLogging(logSetting)
	if err := redis.InitRedisInstance(context.TODO(), cfg.Global.RedisInfo); err != nil {
		logger.Error("[main.go] redis init fail err:%v", err)
		return
	}
	serverInstance := server.GetServer()
	serverInstance.Init(cfg, *pid)
	serverInstance.BaseService.Start()

	logger.Info("[main.go] server close...")
}
