package main

import "greatestworks/world"

func main() {
	world.MM = world.NewMgrMgr()
	go world.MM.Run()
	select {}
}
