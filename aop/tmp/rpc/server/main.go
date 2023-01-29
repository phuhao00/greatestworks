package main

import (
	"context"
	"google.golang.org/grpc"
	r "greatestworks/aop/tmp/rpc"
	"log"
	"net"
)

const (
	address = "localhost:50051"
)

func main() {
	// 监听端口
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer() //获取新服务示例
	r.RegisterGreeterServer(s, &GreeterServer{})
	// 开始处理
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

type GreeterServer struct {
	r.UnimplementedGreeterServer
}

func (g *GreeterServer) Hello(ctx context.Context, request *r.HelloRequest) (*r.HelloResponse, error) {
	resp := &r.HelloResponse{Greeting: "ha"}
	return resp, nil
}
