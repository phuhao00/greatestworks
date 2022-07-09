package main

import "greatestworks/network"

func main() {
	client := network.NewClient(":8023")
	client.Run()
	select {}

}
