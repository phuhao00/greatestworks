package main

import (
	"github.com/phuhao00/network"
)

func main() {
	server := network.NewServer(":8023")
	go server.Run()
	select {}
}
