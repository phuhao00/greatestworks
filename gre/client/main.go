package main

import "github.com/phuhao00/sugar"

func main() {
	c := NewClient()
	c.InputHandlerRegister()
	c.MessageHandlerRegister()
	c.Run()
	sugar.WaitSignal(c.OnSystemSignal)
}
