package main

import (
	"fmt"
	r "greatestworks/aop/tmp/rpc"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

func main() {
	rpc.Register(new(HelloServer)) // 注册rpc服务
	rpc.HandleHTTP()               // 采用http协议作为rpc载体

	lis, err := net.Listen("tcp", "127.0.0.1:8095")
	if err != nil {
		log.Fatalln("fatal error: ", err)
	}

	fmt.Fprintf(os.Stdout, "%s", "start connection")

	http.Serve(lis, nil)

}

type HelloServer struct {
}

func (h *HelloServer) Greet(req r.HelloRequest, res *r.HelloResponse) error {
	res.Greeting = "2222"
	return nil
}
