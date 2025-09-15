# Makefile for MMO Game Server

# 变量定义
APP_NAME := mmo-server
VERSION := 1.0.0
GO_VERSION := 1.21
DOCKER_IMAGE := $(APP_NAME):$(VERSION)
DOCKER_REGISTRY := your-registry.com

# Go相关变量
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := gofmt
GOVET := $(GOCMD) vet
GOLINT := golangci-lint

# 构建目录
BUILD_DIR := ./build
CMD_DIR := ./cmd/server
MAIN_FILE := $(CMD_DIR)/main.go
BINARY_NAME := $(BUILD_DIR)/$(APP_NAME)

# 默认目标
.PHONY: all
all: clean deps fmt vet test build

# 帮助信息
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  clean          - Clean build artifacts"
	@echo "  deps           - Download dependencies"
	@echo "  fmt            - Format Go code"
	@echo "  vet            - Run go vet"
	@echo "  lint           - Run golangci-lint"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  run            - Run the application"
	@echo "  dev            - Run in development mode"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-push    - Push Docker image"
	@echo "  compose-up     - Start with docker-compose"
	@echo "  compose-down   - Stop docker-compose"
	@echo "  deploy         - Deploy to production"
	@echo "  install        - Install the application"
	@echo "  uninstall      - Uninstall the application"

# 创建构建目录
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# 下载依赖
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# 格式化代码
.PHONY: fmt
fmt:
	@echo "Formatting Go code..."
	$(GOFMT) -s -w .

# 代码检查
.PHONY: vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# 代码规范检查
.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	$(GOLINT) run ./...

# 运行测试
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# 运行测试并生成覆盖率报告
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 构建应用
.PHONY: build
build: $(BUILD_DIR)
	@echo "Building $(APP_NAME)..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -a -installsuffix cgo -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Build completed: $(BINARY_NAME)"

# 构建Windows版本
.PHONY: build-windows
build-windows: $(BUILD_DIR)
	@echo "Building $(APP_NAME) for Windows..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -a -installsuffix cgo -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME).exe $(MAIN_FILE)

# 构建macOS版本
.PHONY: build-darwin
build-darwin: $(BUILD_DIR)
	@echo "Building $(APP_NAME) for macOS..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -a -installsuffix cgo -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME)-darwin $(MAIN_FILE)

# 构建所有平台版本
.PHONY: build-all
build-all: build build-windows build-darwin

# 清理构建文件
.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# 运行应用
.PHONY: run
run:
	@echo "Running $(APP_NAME)..."
	$(GOCMD) run $(MAIN_FILE)

# 开发模式运行
.PHONY: dev
dev:
	@echo "Running in development mode..."
	ENV=development $(GOCMD) run $(MAIN_FILE)

# 安装应用
.PHONY: install
install: build
	@echo "Installing $(APP_NAME)..."
	cp $(BINARY_NAME) /usr/local/bin/$(APP_NAME)
	chmod +x /usr/local/bin/$(APP_NAME)

# 卸载应用
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(APP_NAME)..."
	rm -f /usr/local/bin/$(APP_NAME)

# Docker相关命令
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .
	docker tag $(DOCKER_IMAGE) $(APP_NAME):latest

.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run --rm -p 8080:8080 --name $(APP_NAME) $(DOCKER_IMAGE)

.PHONY: docker-push
docker-push: docker-build
	@echo "Pushing Docker image..."
	docker tag $(DOCKER_IMAGE) $(DOCKER_REGISTRY)/$(DOCKER_IMAGE)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE)

# Docker Compose相关命令
.PHONY: compose-up
compose-up:
	@echo "Starting services with docker-compose..."
	docker-compose up -d

.PHONY: compose-down
compose-down:
	@echo "Stopping services with docker-compose..."
	docker-compose down

.PHONY: compose-logs
compose-logs:
	@echo "Showing docker-compose logs..."
	docker-compose logs -f

.PHONY: compose-restart
compose-restart:
	@echo "Restarting services with docker-compose..."
	docker-compose restart

# 开发环境
.PHONY: dev-up
dev-up:
	@echo "Starting development environment..."
	docker-compose --profile development up -d

.PHONY: dev-down
dev-down:
	@echo "Stopping development environment..."
	docker-compose --profile development down

# 生产环境
.PHONY: prod-up
prod-up:
	@echo "Starting production environment..."
	docker-compose --profile production up -d

.PHONY: prod-down
prod-down:
	@echo "Stopping production environment..."
	docker-compose --profile production down

# 监控环境
.PHONY: monitoring-up
monitoring-up:
	@echo "Starting monitoring stack..."
	docker-compose --profile monitoring up -d

.PHONY: monitoring-down
monitoring-down:
	@echo "Stopping monitoring stack..."
	docker-compose --profile monitoring down

# 数据库相关
.PHONY: db-migrate
db-migrate:
	@echo "Running database migrations..."
	# 这里可以添加数据库迁移命令

.PHONY: db-seed
db-seed:
	@echo "Seeding database..."
	# 这里可以添加数据库种子数据命令

.PHONY: db-backup
db-backup:
	@echo "Backing up database..."
	docker-compose exec mongodb mongodump --out /data/backup/$(shell date +%Y%m%d_%H%M%S)

.PHONY: db-restore
db-restore:
	@echo "Restoring database..."
	# docker-compose exec mongodb mongorestore /data/backup/BACKUP_NAME

# 性能测试
.PHONY: benchmark
benchmark:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# 生成文档
.PHONY: docs
docs:
	@echo "Generating documentation..."
	$(GOCMD) doc -all ./... > docs/api.md

# 安全扫描
.PHONY: security
security:
	@echo "Running security scan..."
	gosec ./...

# 依赖检查
.PHONY: deps-check
deps-check:
	@echo "Checking dependencies..."
	$(GOCMD) list -u -m all

# 更新依赖
.PHONY: deps-update
deps-update:
	@echo "Updating dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy

# 完整的CI流程
.PHONY: ci
ci: clean deps fmt vet lint test build

# 部署到生产环境
.PHONY: deploy
deploy: ci docker-build docker-push
	@echo "Deploying to production..."
	# 这里可以添加部署脚本

# 健康检查
.PHONY: health
health:
	@echo "Checking application health..."
	curl -f http://localhost:8080/health || exit 1

# 查看日志
.PHONY: logs
logs:
	@echo "Showing application logs..."
	tail -f logs/app.log

# 监控指标
.PHONY: metrics
metrics:
	@echo "Showing application metrics..."
	curl http://localhost:8080/metrics