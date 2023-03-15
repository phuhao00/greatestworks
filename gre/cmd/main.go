package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	usage = "hhhhhhhhh"
)

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}
	println("kkkkk")

}
