#!/bin/bash
# 健康检查脚本

set -e

echo "检查服务健康状态..."

# 检查HTTP服务器
echo "检查HTTP服务器..."
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "✓ HTTP服务器正常"
else
    echo "✗ HTTP服务器异常"
    exit 1
fi

# 检查MongoDB
echo "检查MongoDB..."
if docker exec mmo-server mongosh --host mongodb:27017 --username admin --password admin123 --authenticationDatabase admin --eval "db.runCommand('ping')" > /dev/null 2>&1; then
    echo "✓ MongoDB正常"
else
    echo "✗ MongoDB异常"
    exit 1
fi

# 检查Redis
echo "检查Redis..."
if docker exec mmo-server redis-cli -h redis -p 6379 -a redis123 ping > /dev/null 2>&1; then
    echo "✓ Redis正常"
else
    echo "✗ Redis异常"
    exit 1
fi

# 检查NATS
echo "检查NATS..."
if nc -z localhost 4222; then
    echo "✓ NATS正常"
else
    echo "✗ NATS异常"
    exit 1
fi

echo "所有服务健康检查通过！"


