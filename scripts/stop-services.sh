#!/bin/bash
# 停止服务脚本

set -e

echo "停止GreatestWorks MMO游戏服务器..."

# 停止所有服务
echo "停止Docker Compose服务..."
docker-compose down

# 清理数据卷（可选）
read -p "是否清理数据卷？这将删除所有数据 (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "清理数据卷..."
    docker-compose down -v
    docker volume prune -f
fi

# 清理镜像（可选）
read -p "是否清理构建镜像？ (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "清理构建镜像..."
    docker image prune -f
fi

echo "服务已停止"
