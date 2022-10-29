package main

import (
	"github.com/phuhao00/sugar"
	"greatestworks/aop/logger"
	"greatestworks/business/server/world"
)

func main() {
	world.Oasis = world.NewWorld()
	go world.Oasis.Start()
	logger.Logger.InfoF("server start !!")
	sugar.WaitSignal(world.Oasis.OnSystemSignal)
}
