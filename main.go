package main

import (
	"greatestworks/logger"
	"greatestworks/world"

	"github.com/phuhao00/sugar"
)

func main() {
	world.MM = world.NewMgrMgr()
	go world.MM.Run()
	logger.Logger.InfoF("server start !!")
	sugar.WaitSignal(world.MM.OnSystemSignal)
}
