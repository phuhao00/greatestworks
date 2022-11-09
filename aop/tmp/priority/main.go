package main

import "fmt"

var name = "huhao"

func init() {
	fmt.Println(name)
	name = "huhao1"

}

func main() {
	fmt.Println(name)

}
