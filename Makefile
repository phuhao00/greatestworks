# GreatestWorks MMO 游戏服务器 Makefile

.PHONY: help build run stop clean test lint format docker-build docker-run docker-stop health-check

# 默认目标
help:
	@echo "GreatestWorks MMO 游戏服务器"
	@echo "可用命令:"
	@echo "  build          - 构建项目"
	@echo "  run            - 运行服务"
	@echo "  stop           - 停止服务"
	@echo "  clean          - 清理构建文件"
	@echo "  test           - 运行测试"
	@echo "  lint           - 代码检查"
	@echo "  format         - 格式化代码"
	@echo "  docker-build   - 构建Docker镜像"
	@echo "  docker-run     - 运行Docker服务"
	@echo "  docker-stop    - 停止Docker服务"
	@echo "  health-check   - 健康检查"

# 构建项目
build:
	@echo "构建项目..."
	go build -o bin/game-service ./cmd/game-service
	go build -o bin/auth-service ./cmd/auth-service
	go build -o bin/gateway-service ./cmd/gateway-service

# 运行服务
run: build
	@echo "启动游戏服务..."
	./bin/game-service

# 停止服务
stop:
	@echo "停止服务..."
	pkill -f game-service || true

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf bin/
	rm -rf logs/
	docker system prune -f

# 运行测试
test:
	@echo "运行测试..."
	go test -v ./...

# 代码检查
lint:
	@echo "代码检查..."
	golangci-lint run

# 格式化代码
format:
	@echo "格式化代码..."
	go fmt ./...
	goimports -w .

# 构建Docker镜像
docker-build:
	@echo "构建Docker镜像..."
	docker build -t greatestworks/mmo-server:latest .

# 运行Docker服务
docker-run:
	@echo "启动Docker服务..."
	chmod +x scripts/*.sh
	docker-compose up -d

# 停止Docker服务
docker-stop:
	@echo "停止Docker服务..."
	docker-compose down

# 健康检查
health-check:
	@echo "健康检查..."
	chmod +x scripts/health-check.sh
	./scripts/health-check.sh

# 开发环境快速启动
dev: docker-run
	@echo "开发环境已启动"
	@echo "HTTP服务器: http://localhost:8080"
	@echo "健康检查: http://localhost:8080/health"
	@echo "指标监控: http://localhost:8080/metrics"

# 生产环境部署
deploy: docker-build
	@echo "部署到生产环境..."
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# 查看日志
logs:
	@echo "查看服务日志..."
	docker-compose logs -f

# 进入容器
shell:
	@echo "进入游戏服务容器..."
	docker exec -it mmo-server /bin/sh

# 运行模拟客户端
simclient:
	@echo "运行模拟客户端 (integration 模式)..."
	go run ./tools/simclient/cmd/simclient -mode integration