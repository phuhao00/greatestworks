package main

import (
	"github.com/phuhao00/sugar"
	"greatestworks/aop/logger"
	"greatestworks/business/server/world"
)

func main() {
	world.MM = world.NewMgrMgr()
	go world.MM.Start()
	logger.Logger.InfoF("server start !!")
	sugar.WaitSignal(world.MM.OnSystemSignal)
}
