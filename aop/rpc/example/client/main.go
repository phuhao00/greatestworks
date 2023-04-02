package main

import (
	"fmt"
	"greatestworks/aop/rpc"
	"greatestworks/aop/rpc/example/processor"
)

func main() {
	c := rpc.NewRpcClient("127.0.0.1:444")
	req := &processor.MockParam{Tag: "123"}
	rsp := &processor.MockParam{Tag: ""}
	c.Call("MockProcessor.Print2", req, rsp)
	fmt.Println(rsp)

}
