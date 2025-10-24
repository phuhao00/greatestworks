#!/bin/bash
# 数据库初始化脚本

set -e

echo "开始初始化数据库..."

# 等待MongoDB启动
echo "等待MongoDB启动..."
until mongosh --host mongodb:27017 --eval "print('MongoDB is ready')" > /dev/null 2>&1; do
  echo "等待MongoDB启动..."
  sleep 2
done

# 等待Redis启动
echo "等待Redis启动..."
until redis-cli -h redis -p 6379 -a redis123 ping > /dev/null 2>&1; do
  echo "等待Redis启动..."
  sleep 2
done

# 等待NATS启动
echo "等待NATS启动..."
until nc -z nats 4222; do
  echo "等待NATS启动..."
  sleep 2
done

echo "所有数据库服务已启动"

# 执行MongoDB初始化脚本
echo "执行MongoDB初始化..."
mongosh --host mongodb:27017 --username admin --password admin123 --authenticationDatabase admin < /scripts/mongo-init.js

# 初始化Redis数据
echo "初始化Redis数据..."
redis-cli -h redis -p 6379 -a redis123 << EOF
FLUSHDB
SET server:status "running"
SET server:start_time "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
SET players:count 0
SET battles:count 0
SET items:count 0
EOF

echo "数据库初始化完成"


