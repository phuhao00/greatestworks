#!/bin/bash

# Proto文件生成脚本
# 支持Go和C#代码生成，只生成消息定义，不生成gRPC服务

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}开始生成Proto文件...${NC}"

# 检查protoc是否安装
if ! command -v protoc &> /dev/null; then
    echo -e "${RED}错误: protoc未安装，请先安装Protocol Buffers编译器${NC}"
    echo "安装方法:"
    echo "  macOS: brew install protobuf"
    echo "  Ubuntu: sudo apt-get install protobuf-compiler"
    echo "  CentOS: sudo yum install protobuf-compiler"
    exit 1
fi

# 检查Go插件
if ! command -v protoc-gen-go &> /dev/null; then
    echo -e "${YELLOW}安装Go插件...${NC}"
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# 检查C#插件
if ! command -v protoc-gen-csharp &> /dev/null; then
    echo -e "${YELLOW}安装C#插件...${NC}"
    # C#插件通常随protoc一起安装
    echo "请确保已安装C#支持"
fi

# 创建输出目录
mkdir -p internal/proto/player
mkdir -p internal/proto/battle
mkdir -p internal/proto/pet
mkdir -p internal/proto/common
mkdir -p csharp/GreatestWorks/Player
mkdir -p csharp/GreatestWorks/Battle
mkdir -p csharp/GreatestWorks/Pet
mkdir -p csharp/GreatestWorks/Common

echo -e "${GREEN}生成Go代码（只生成消息定义，不生成gRPC服务）...${NC}"

# 生成Go代码 - 只生成消息定义，不生成gRPC服务
protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    proto/common.proto

protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    proto/player.proto

protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    proto/battle.proto

protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    proto/pet.proto

echo -e "${GREEN}移动生成的文件到正确位置...${NC}"

# 移动生成的文件到正确位置
mv proto/common.pb.go internal/proto/common/ 2>/dev/null || true
mv proto/player.pb.go internal/proto/player/ 2>/dev/null || true
mv proto/battle.pb.go internal/proto/battle/ 2>/dev/null || true
mv proto/pet.pb.go internal/proto/pet/ 2>/dev/null || true

# 删除可能生成的gRPC文件
rm -f proto/*_grpc.pb.go

echo -e "${GREEN}生成C#代码...${NC}"

# 生成C#代码
protoc \
    --csharp_out=csharp \
    --csharp_opt=file_extension=.g.cs \
    proto/common.proto

protoc \
    --csharp_out=csharp \
    --csharp_opt=file_extension=.g.cs \
    proto/player.proto

protoc \
    --csharp_out=csharp \
    --csharp_opt=file_extension=.g.cs \
    proto/battle.proto

protoc \
    --csharp_out=csharp \
    --csharp_opt=file_extension=.g.cs \
    proto/pet.proto

echo -e "${GREEN}Proto文件生成完成！${NC}"
echo -e "${YELLOW}生成的文件:${NC}"
echo "  Go代码: internal/proto/"
echo "  C#代码: csharp/"

# 显示生成的文件
echo -e "${YELLOW}Go生成的文件:${NC}"
find internal/proto -name "*.pb.go" -type f

echo -e "${YELLOW}C#生成的文件:${NC}"
find csharp -name "*.g.cs" -type f

echo -e "${YELLOW}注意: 只生成了protobuf消息定义，没有生成gRPC服务文件${NC}"
echo -e "${YELLOW}项目使用netcore-go RPC架构，不使用gRPC${NC}"
