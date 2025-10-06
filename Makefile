# GreatestWorks Makefile
# 支持Proto文件生成和项目构建

.PHONY: help proto-go proto-csharp proto-all clean build test

# 默认目标
help:
	@echo "GreatestWorks 构建系统"
	@echo ""
	@echo "可用命令:"
	@echo "  proto-go     生成Go的Proto文件"
	@echo "  proto-csharp 生成C#的Proto文件"
	@echo "  proto-all    生成所有Proto文件"
	@echo "  build        构建项目"
	@echo "  test         运行测试"
	@echo "  clean        清理生成的文件"
	@echo "  install-deps 安装依赖"

# 检查protoc是否安装
check-protoc:
	@which protoc > /dev/null || (echo "错误: protoc未安装" && exit 1)

# 检查Go插件
check-go-plugins:
	@which protoc-gen-go > /dev/null || (echo "安装Go插件..." && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest)
	@which protoc-gen-go-grpc > /dev/null || (echo "安装Go gRPC插件..." && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest)

# 生成Go Proto文件
proto-go: check-protoc check-go-plugins
	@echo "生成Go Proto文件（只生成消息定义，不生成gRPC服务）..."
	@mkdir -p internal/proto/player
	@mkdir -p internal/proto/battle
	@mkdir -p internal/proto/pet
	@mkdir -p internal/proto/common
	@protoc --go_out=. --go_opt=paths=source_relative proto/common.proto
	@protoc --go_out=. --go_opt=paths=source_relative proto/player.proto
	@protoc --go_out=. --go_opt=paths=source_relative proto/battle.proto
	@protoc --go_out=. --go_opt=paths=source_relative proto/pet.proto
	@echo "移动生成的文件到正确位置..."
	@mv proto/common.pb.go internal/proto/common/ 2>/dev/null || true
	@mv proto/player.pb.go internal/proto/player/ 2>/dev/null || true
	@mv proto/battle.pb.go internal/proto/battle/ 2>/dev/null || true
	@mv proto/pet.pb.go internal/proto/pet/ 2>/dev/null || true
	@rm -f proto/*_grpc.pb.go
	@echo "Go Proto文件生成完成"

# 生成C# Proto文件
proto-csharp: check-protoc
	@echo "生成C# Proto文件..."
	@mkdir -p csharp/GreatestWorks/Player
	@mkdir -p csharp/GreatestWorks/Battle
	@mkdir -p csharp/GreatestWorks/Pet
	@mkdir -p csharp/GreatestWorks/Common
	@protoc --csharp_out=csharp --csharp_opt=file_extension=.g.cs proto/common.proto
	@protoc --csharp_out=csharp --csharp_opt=file_extension=.g.cs proto/player.proto
	@protoc --csharp_out=csharp --csharp_opt=file_extension=.g.cs proto/battle.proto
	@protoc --csharp_out=csharp --csharp_opt=file_extension=.g.cs proto/pet.proto
	@echo "C# Proto文件生成完成"

# 生成所有Proto文件
proto-all: proto-go proto-csharp
	@echo "所有Proto文件生成完成"

# 构建项目
build: proto-go
	@echo "构建项目..."
	@go mod tidy
	@go build -o bin/server cmd/server/main.go
	@echo "构建完成: bin/server"

# 运行测试
test:
	@echo "运行测试..."
	@go test ./...

# 清理生成的文件
clean:
	@echo "清理生成的文件..."
	@rm -rf internal/proto/*/
	@rm -rf csharp/
	@rm -rf bin/
	@go clean
	@echo "清理完成"

# 安装依赖
install-deps:
	@echo "安装Go依赖..."
	@go mod download
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "依赖安装完成"

# 开发环境设置
dev-setup: install-deps proto-all
	@echo "开发环境设置完成"

# 运行服务器
run: build
	@echo "启动服务器..."
	@./bin/server

# 格式化代码
fmt:
	@echo "格式化代码..."
	@go fmt ./...
	@echo "代码格式化完成"

# 检查代码
lint:
	@echo "检查代码..."
	@go vet ./...
	@echo "代码检查完成"