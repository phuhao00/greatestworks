#!/bin/bash
# 启动服务脚本

set -e

echo "启动GreatestWorks MMO游戏服务器..."

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo "错误: Docker未运行，请先启动Docker"
    exit 1
fi

# 检查Docker Compose是否可用
if ! command -v docker-compose > /dev/null 2>&1; then
    echo "错误: Docker Compose未安装"
    exit 1
fi

# 创建必要的目录
mkdir -p logs
mkdir -p configs
mkdir -p data/mongodb
mkdir -p data/redis
mkdir -p data/nats

# 设置权限
chmod +x scripts/*.sh

# 复制环境变量文件（如果不存在）
if [ ! -f .env ]; then
    echo "创建环境变量文件..."
    cat > .env << EOF
# 构建配置
BUILD_TARGET=final
BUILD_VERSION=1.0.0
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT=dev
IMAGE_TAG=latest

# 应用配置
APP_ENV=development
GIN_MODE=debug
LOG_LEVEL=info
LOG_FORMAT=json

# 服务器端口配置
SERVER_HTTP_PORT=8080
SERVER_WS_PORT=8081
SERVER_METRICS_PORT=9090

# 数据库配置
MONGODB_USER=admin
MONGODB_PASSWORD=admin123
MONGODB_DATABASE=mmo_game
REDIS_PASSWORD=redis123

# 消息队列配置
NATS_CLUSTER_ID=mmo-cluster

# 安全配置
JWT_SECRET=your-super-secret-jwt-key-change-in-production
ENCRYPTION_KEY=32-character-encryption-key-here

# 性能配置
MAX_CONNECTIONS=10000
WORKER_POOL_SIZE=100
CACHE_TTL=3600

# 资源限制
SERVER_CPU_LIMIT=2.0
SERVER_MEMORY_LIMIT=2G
SERVER_CPU_RESERVATION=0.5
SERVER_MEMORY_RESERVATION=512M
EOF
fi

# 启动服务
echo "启动Docker Compose服务..."
docker-compose up -d

# 等待服务启动
echo "等待服务启动..."
sleep 10

# 检查服务状态
echo "检查服务状态..."
docker-compose ps

# 显示日志
echo "显示服务日志..."
docker-compose logs --tail=50

echo "服务启动完成！"
echo "HTTP服务器: http://localhost:8080"
echo "健康检查: http://localhost:8080/health"
echo "指标监控: http://localhost:8080/metrics"
echo "MongoDB管理: http://localhost:8081"
echo "Redis管理: http://localhost:8082"