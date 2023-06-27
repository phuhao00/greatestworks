package server

import (
	"greatestworks/aop/logger"
)

func (w *World) Reload() {
	logger.Info("[Reload] World Reload ")
}

func (w *World) Init(config interface{}, processId int) {

}

func (w *World) Start() {
	w.HandlerRegister()
	go w.Run()

}

func (w *World) Stop() {

}
