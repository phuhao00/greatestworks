package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	r "greatestworks/aop/tmp/rpc"
	"log"
)

const (
	address = "localhost:50051"
)

var cc *grpc.ClientConn
var gc r.GreeterClient

func main() {
	//建立链接
	var err error
	cc, err = grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	gc = r.NewGreeterClient(cc)
	defer cc.Close()
	response, err := gc.Hello(context.TODO(), &r.HelloRequest{Name: "hello"})
	if err != nil {

	}
	fmt.Println(response)
}
