package main

import "greatestworks/world"

func main() {
	world.MM = world.NewMgrMgr()
	world.MM.Run()
}
