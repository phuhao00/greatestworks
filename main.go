package main

import (
	"greatestworks/business/world"
	"greatestworks/logger"

	"github.com/phuhao00/sugar"
)

func main() {
	world.MM = world.NewMgrMgr()
	go world.MM.Run()
	logger.Logger.InfoF("server start !!")
	sugar.WaitSignal(world.MM.OnSystemSignal)
}
