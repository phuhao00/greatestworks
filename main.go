package main

import (
	"greatestworks/aop/logger"
	"greatestworks/business/world"

	"github.com/phuhao00/sugar"
)

func main() {
	world.MM = world.NewMgrMgr()
	go world.MM.Start()
	logger.Logger.InfoF("server start !!")
	sugar.WaitSignal(world.MM.OnSystemSignal)
}
