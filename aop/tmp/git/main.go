package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
)

func main() {
	repository, err := git.PlainOpen(".")
	if err != nil {

	}
	head, err := repository.Head()
	if err != nil {

	}
	fmt.Println(head.Name())
}
