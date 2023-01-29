package main

import (
	"fmt"
	r "greatestworks/aop/tmp/rpc"
	"log"
	"net/rpc"
)

func main() {
	conn, err := rpc.DialHTTP("tcp", "127.0.0.1:8095")
	if err != nil {
		log.Fatalln("dailing error: ", err)
	}

	req := r.HelloRequest{Name: "111"}
	var res r.HelloResponse
	err = conn.Call("HelloServer.Greet", req, &res) // 乘法运算
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Printf("%s\n", res.Greeting)
}
