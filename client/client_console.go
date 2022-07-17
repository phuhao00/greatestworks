package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ClientConsole struct {
	chInput chan *InputParam
}

type InputParam struct {
	Command string
	Param   []string
}

func NewClientConsole() *ClientConsole {
	c := &ClientConsole{}
	return c
}

func (c *ClientConsole) Run() {
	reader := bufio.NewReader(os.Stdin)
	for {
		readString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("input err ,check your input and  try again !!!")
			continue
		}
		split := strings.Split(readString, "|")
		newSlice := make([]string, 0)
		for i, s := range split {
			split[i] = strings.TrimSpace(s)
			split[i] = strings.Trim(s, "\n")
			split[i] = strings.Trim(s, "\r")
			if len(s) != 0 {
				newSlice = append(newSlice, s)
			}
		}
		if len(newSlice) == 0 {
			fmt.Println("input err, check your input and  try again !!! ")
			continue
		}
		in := &InputParam{
			Command: newSlice[0],
			Param:   newSlice[1:],
		}
		c.chInput <- in
	}
}
