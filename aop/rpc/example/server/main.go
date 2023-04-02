package main

import (
	"fmt"
	"github.com/phuhao00/spoor"
	"greatestworks/aop/logger"
	rpcLocal "greatestworks/aop/rpc"
	"greatestworks/aop/rpc/example/processor"
	"net/rpc"
)

func main() {
	logger.SetLogging(&logger.LoggingSetting{
		Dir:          "./log",
		Level:        int(spoor.DEBUG),
		Prefix:       "",
		WriterOption: nil,
	})
	if err := rpc.Register(new(processor.MockProcessor)); err != nil {
		fmt.Println("register failed")
		return
	}
	rpcServer := &rpcLocal.Server{}
	rpcServer.Init("127.0.0.1:444")

	rpcServer.Run()
}
