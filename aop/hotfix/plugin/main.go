package main

import (
	"fmt"
	"os"
	"plugin"
	"time"
)

func main() {
	//todo
	time.Sleep(time.Second * 2)
	Hello()
	time.Sleep(time.Second * 2)
}

func Hello() {
	p, err := plugin.Open("./plugina.so")
	if err != nil {
		fmt.Println("error open plugin: ", err)
		os.Exit(-1)
	}
	s, err := p.Lookup("IamPluginA")
	if err != nil {
		fmt.Println("error lookup IamPluginA: ", err)
		os.Exit(-1)
	}
	if x, ok := s.(func()); ok {
		x()
	}
}
