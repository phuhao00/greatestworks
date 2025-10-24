#!/bin/bash

# Greatest Works - 构建脚本
# 用于编译Go项目，支持多平台交叉编译

set -e

# 项目信息
PROJECT_NAME="greatestworks"
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 构建参数
LDFLAGS="-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}"
BUILD_DIR="./bin"
SOURCE_DIR="./cmd/server"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 清理构建目录
clean_build_dir() {
    log_info "清理构建目录..."
    rm -rf ${BUILD_DIR}
    mkdir -p ${BUILD_DIR}
}

# 检查Go环境
check_go_env() {
    if ! command -v go &> /dev/null; then
        log_error "Go未安装或不在PATH中"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    log_info "Go版本: ${GO_VERSION}"
    
    # 检查Go版本是否满足要求
    REQUIRED_VERSION="1.21"
    if ! printf '%s\n%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V -C; then
        log_warning "建议使用Go ${REQUIRED_VERSION}或更高版本"
    fi
}

# 下载依赖
download_deps() {
    log_info "下载Go模块依赖..."
    go mod download
    go mod tidy
}

# 运行测试
run_tests() {
    if [ "$SKIP_TESTS" != "true" ]; then
        log_info "运行单元测试..."
        go test -v ./... -timeout=30s
        if [ $? -ne 0 ]; then
            log_error "测试失败，构建中止"
            exit 1
        fi
        log_success "所有测试通过"
    else
        log_warning "跳过测试"
    fi
}

# 代码质量检查
run_lint() {
    if [ "$SKIP_LINT" != "true" ]; then
        if command -v golangci-lint &> /dev/null; then
            log_info "运行代码质量检查..."
            golangci-lint run
            if [ $? -ne 0 ]; then
                log_error "代码质量检查失败，构建中止"
                exit 1
            fi
            log_success "代码质量检查通过"
        else
            log_warning "golangci-lint未安装，跳过代码质量检查"
        fi
    else
        log_warning "跳过代码质量检查"
    fi
}

# 构建单个平台
build_single() {
    local goos=$1
    local goarch=$2
    local output_name="${PROJECT_NAME}"
    
    if [ "$goos" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    local output_path="${BUILD_DIR}/${goos}-${goarch}/${output_name}"
    
    log_info "构建 ${goos}/${goarch}..."
    
    mkdir -p "$(dirname "$output_path")"
    
    CGO_ENABLED=0 GOOS=$goos GOARCH=$goarch go build \
        -ldflags "${LDFLAGS}" \
        -o "$output_path" \
        "$SOURCE_DIR"
    
    if [ $? -eq 0 ]; then
        local file_size=$(du -h "$output_path" | cut -f1)
        log_success "构建完成: $output_path (${file_size})"
    else
        log_error "构建失败: ${goos}/${goarch}"
        return 1
    fi
}

# 构建所有平台
build_all() {
    log_info "开始多平台构建..."
    
    # 定义支持的平台
    declare -a platforms=(
        "linux/amd64"
        "linux/arm64"
        "darwin/amd64"
        "darwin/arm64"
        "windows/amd64"
    )
    
    for platform in "${platforms[@]}"; do
        IFS='/' read -r goos goarch <<< "$platform"
        build_single "$goos" "$goarch"
        if [ $? -ne 0 ]; then
            log_error "构建失败，中止"
            exit 1
        fi
    done
    
    log_success "所有平台构建完成"
}

# 构建当前平台
build_current() {
    local goos=$(go env GOOS)
    local goarch=$(go env GOARCH)
    
    log_info "构建当前平台 ${goos}/${goarch}..."
    
    local output_name="${PROJECT_NAME}"
    if [ "$goos" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    local output_path="${BUILD_DIR}/${output_name}"
    
    go build -ldflags "${LDFLAGS}" -o "$output_path" "$SOURCE_DIR"
    
    if [ $? -eq 0 ]; then
        local file_size=$(du -h "$output_path" | cut -f1)
        log_success "构建完成: $output_path (${file_size})"
        
        # 创建符号链接到项目根目录
        ln -sf "$output_path" "./server"
        log_info "创建符号链接: ./server -> $output_path"
    else
        log_error "构建失败"
        exit 1
    fi
}

# 显示帮助信息
show_help() {
    echo "Greatest Works 构建脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help     显示帮助信息"
    echo "  -a, --all      构建所有支持的平台"
    echo "  -c, --current  构建当前平台 (默认)"
    echo "  --skip-tests   跳过单元测试"
    echo "  --skip-lint    跳过代码质量检查"
    echo "  --clean        清理构建目录"
    echo ""
    echo "环境变量:"
    echo "  VERSION        版本号 (默认: git describe)"
    echo "  SKIP_TESTS     跳过测试 (true/false)"
    echo "  SKIP_LINT      跳过代码检查 (true/false)"
    echo ""
    echo "示例:"
    echo "  $0                    # 构建当前平台"
    echo "  $0 --all              # 构建所有平台"
    echo "  $0 --skip-tests       # 跳过测试构建"
    echo "  VERSION=v1.0.0 $0     # 指定版本号构建"
}

# 主函数
main() {
    local build_all_platforms=false
    local clean_only=false
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -a|--all)
                build_all_platforms=true
                shift
                ;;
            -c|--current)
                build_all_platforms=false
                shift
                ;;
            --skip-tests)
                export SKIP_TESTS=true
                shift
                ;;
            --skip-lint)
                export SKIP_LINT=true
                shift
                ;;
            --clean)
                clean_only=true
                shift
                ;;
            *)
                log_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    log_info "Greatest Works 构建脚本启动"
    log_info "版本: ${VERSION}"
    log_info "构建时间: ${BUILD_TIME}"
    log_info "Git提交: ${GIT_COMMIT}"
    
    # 清理构建目录
    clean_build_dir
    
    if [ "$clean_only" = true ]; then
        log_success "构建目录已清理"
        exit 0
    fi
    
    # 检查环境
    check_go_env
    
    # 下载依赖
    download_deps
    
    # 运行测试
    run_tests
    
    # 代码质量检查
    run_lint
    
    # 构建
    if [ "$build_all_platforms" = true ]; then
        build_all
    else
        build_current
    fi
    
    log_success "构建完成！"
}

# 执行主函数
main "$@"