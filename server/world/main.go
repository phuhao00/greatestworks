package main

import (
	"github.com/phuhao00/sugar"
	"greatestworks/aop/logger"
)

func main() {
	Oasis = NewWorld()
	go Oasis.Start()
	logger.InfoF("server start !!")
	sugar.WaitSignal(Oasis.OnSystemSignal)
}
