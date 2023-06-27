package main

import (
	"github.com/phuhao00/sugar"
	"greatestworks/aop/logger"
	"greatestworks/server/world/server"
)

func main() {
	server.Oasis = server.NewWorld()
	go server.Oasis.Start()
	logger.Info("server start !!")
	sugar.WaitSignal(server.Oasis.OnSystemSignal)
}
