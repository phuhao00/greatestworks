// Package main 游戏服务器主程序
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 创建服务器启动器
	bootstrap := NewServerBootstrap()
	
	// 初始化服务器
	if err := bootstrap.Initialize(); err != nil {
		log.Fatalf("服务器初始化失败: %v", err)
	}
	
	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	bootstrap.bootstrap.ShutdownChan = sigChan
	
	// 启动服务器
	if err := bootstrap.StartServer(); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}