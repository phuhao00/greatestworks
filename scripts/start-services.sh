#!/bin/bash
# 启动分布式游戏服务脚本
# 基于DDD架构的分布式多节点服务

echo "启动分布式游戏服务..."

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "错误: 未找到Go环境，请先安装Go"
    exit 1
fi

# 创建日志目录
mkdir -p logs

# 启动认证服务
echo "启动认证服务..."
gnome-terminal --title="Auth Service" -- bash -c "cd $(dirname $0)/.. && go run cmd/auth-service/main.go; exec bash" &

# 等待认证服务启动
sleep 3

# 启动游戏服务
echo "启动游戏服务..."
gnome-terminal --title="Game Service" -- bash -c "cd $(dirname $0)/.. && go run cmd/game-service/main.go; exec bash" &

# 等待游戏服务启动
sleep 3

# 启动网关服务
echo "启动网关服务..."
gnome-terminal --title="Gateway Service" -- bash -c "cd $(dirname $0)/.. && go run cmd/gateway-service/main.go; exec bash" &

echo ""
echo "所有服务已启动！"
echo ""
echo "服务地址："
echo "- 认证服务: http://localhost:8080"
echo "- 游戏服务: rpc://localhost:8081"
echo "- 网关服务: tcp://localhost:9090"
echo ""
echo "按Ctrl+C关闭所有服务..."

# 等待用户中断
trap 'echo "正在关闭所有服务..."; pkill -f "go run"; exit 0' INT

# 保持脚本运行
while true; do
    sleep 1
done
