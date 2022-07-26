package main

import (
	"github.com/phuhao00/network"
)

func main() {
	client := network.NewClient(":8023")
	client.Run()
	select {}

}
