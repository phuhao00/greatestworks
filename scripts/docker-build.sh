#!/bin/bash

# Greatest Works - Docker构建和推送脚本
# 支持多架构构建和镜像推送

set -e

# 默认配置
DEFAULT_REGISTRY="registry.example.com"
DEFAULT_NAMESPACE="greatestworks"
PROJECT_NAME="greatestworks"
DEFAULT_TAG="latest"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

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

# 检查Docker环境
check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker未安装"
        exit 1
    fi
    
    if ! docker info >/dev/null 2>&1; then
        log_error "Docker未运行，请启动Docker服务"
        exit 1
    fi
    
    log_success "Docker环境检查通过"
}

# 检查buildx支持
check_buildx() {
    if ! docker buildx version >/dev/null 2>&1; then
        log_error "Docker Buildx未安装，无法进行多架构构建"
        return 1
    fi
    
    # 创建buildx构建器（如果不存在）
    if ! docker buildx ls | grep -q "multiarch"; then
        log_info "创建多架构构建器..."
        docker buildx create --name multiarch --driver docker-container --use
        docker buildx inspect --bootstrap
    else
        docker buildx use multiarch
    fi
    
    log_success "Docker Buildx准备完成"
}

# 获取版本信息
get_version_info() {
    # 从git获取版本信息
    if git rev-parse --git-dir >/dev/null 2>&1; then
        VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
        GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
        GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
    else
        VERSION="dev"
        GIT_COMMIT="unknown"
        GIT_BRANCH="unknown"
    fi
    
    BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
    
    log_info "版本信息:"
    log_info "  版本: $VERSION"
    log_info "  提交: $GIT_COMMIT"
    log_info "  分支: $GIT_BRANCH"
    log_info "  构建时间: $BUILD_TIME"
}

# 构建单架构镜像
build_single_arch() {
    local registry="$1"
    local namespace="$2"
    local tag="$3"
    local platform="$4"
    local push="$5"
    
    local image_name="$registry/$namespace/$PROJECT_NAME:$tag"
    
    log_info "构建镜像: $image_name (平台: $platform)"
    
    local build_args=()
    build_args+=("--platform=$platform")
    build_args+=("--build-arg=VERSION=$VERSION")
    build_args+=("--build-arg=GIT_COMMIT=$GIT_COMMIT")
    build_args+=("--build-arg=GIT_BRANCH=$GIT_BRANCH")
    build_args+=("--build-arg=BUILD_TIME=$BUILD_TIME")
    build_args+=("-t=$image_name")
    
    if [ "$push" = true ]; then
        build_args+=("--push")
    else
        build_args+=("--load")
    fi
    
    build_args+=(".")
    
    if docker buildx build "${build_args[@]}"; then
        log_success "镜像构建成功: $image_name"
        return 0
    else
        log_error "镜像构建失败: $image_name"
        return 1
    fi
}

# 构建多架构镜像
build_multi_arch() {
    local registry="$1"
    local namespace="$2"
    local tag="$3"
    local platforms="$4"
    local push="$5"
    
    local image_name="$registry/$namespace/$PROJECT_NAME:$tag"
    
    log_info "构建多架构镜像: $image_name"
    log_info "支持平台: $platforms"
    
    local build_args=()
    build_args+=("--platform=$platforms")
    build_args+=("--build-arg=VERSION=$VERSION")
    build_args+=("--build-arg=GIT_COMMIT=$GIT_COMMIT")
    build_args+=("--build-arg=GIT_BRANCH=$GIT_BRANCH")
    build_args+=("--build-arg=BUILD_TIME=$BUILD_TIME")
    build_args+=("-t=$image_name")
    
    if [ "$push" = true ]; then
        build_args+=("--push")
    fi
    
    build_args+=(".")
    
    if docker buildx build "${build_args[@]}"; then
        log_success "多架构镜像构建成功: $image_name"
        return 0
    else
        log_error "多架构镜像构建失败: $image_name"
        return 1
    fi
}

# 构建开发镜像
build_dev_image() {
    local tag="dev"
    
    log_info "构建开发镜像..."
    
    # 使用开发Dockerfile（如果存在）
    local dockerfile="Dockerfile"
    if [ -f "Dockerfile.dev" ]; then
        dockerfile="Dockerfile.dev"
    fi
    
    local image_name="$PROJECT_NAME:$tag"
    
    docker build \
        -f "$dockerfile" \
        --build-arg VERSION="$VERSION" \
        --build-arg GIT_COMMIT="$GIT_COMMIT" \
        --build-arg BUILD_TIME="$BUILD_TIME" \
        -t "$image_name" \
        .
    
    if [ $? -eq 0 ]; then
        log_success "开发镜像构建成功: $image_name"
        return 0
    else
        log_error "开发镜像构建失败"
        return 1
    fi
}

# 登录到镜像仓库
login_registry() {
    local registry="$1"
    
    if [ -n "$REGISTRY_USERNAME" ] && [ -n "$REGISTRY_PASSWORD" ]; then
        log_info "登录到镜像仓库: $registry"
        echo "$REGISTRY_PASSWORD" | docker login "$registry" -u "$REGISTRY_USERNAME" --password-stdin
        
        if [ $? -eq 0 ]; then
            log_success "镜像仓库登录成功"
            return 0
        else
            log_error "镜像仓库登录失败"
            return 1
        fi
    else
        log_warning "未设置镜像仓库凭据，跳过登录"
        return 0
    fi
}

# 推送镜像
push_image() {
    local image_name="$1"
    
    log_info "推送镜像: $image_name"
    
    if docker push "$image_name"; then
        log_success "镜像推送成功: $image_name"
        return 0
    else
        log_error "镜像推送失败: $image_name"
        return 1
    fi
}

# 标记镜像
tag_image() {
    local source_tag="$1"
    local target_tag="$2"
    
    log_info "标记镜像: $source_tag -> $target_tag"
    
    if docker tag "$source_tag" "$target_tag"; then
        log_success "镜像标记成功"
        return 0
    else
        log_error "镜像标记失败"
        return 1
    fi
}

# 清理构建缓存
clean_build_cache() {
    log_info "清理Docker构建缓存..."
    
    # 清理buildx缓存
    docker buildx prune -f
    
    # 清理未使用的镜像
    docker image prune -f
    
    # 清理悬空镜像
    docker image prune -a -f --filter "until=24h"
    
    log_success "构建缓存清理完成"
}

# 显示镜像信息
show_image_info() {
    local image_name="$1"
    
    log_info "镜像信息: $image_name"
    
    if docker image inspect "$image_name" >/dev/null 2>&1; then
        local size=$(docker image inspect "$image_name" --format='{{.Size}}' | numfmt --to=iec)
        local created=$(docker image inspect "$image_name" --format='{{.Created}}')
        local arch=$(docker image inspect "$image_name" --format='{{.Architecture}}')
        
        echo "  大小: $size"
        echo "  创建时间: $created"
        echo "  架构: $arch"
        
        # 显示层信息
        echo "  层数: $(docker history "$image_name" --quiet | wc -l)"
    else
        log_warning "镜像不存在: $image_name"
    fi
}

# 运行安全扫描
run_security_scan() {
    local image_name="$1"
    
    log_info "运行安全扫描: $image_name"
    
    # 检查是否安装了安全扫描工具
    if command -v trivy >/dev/null 2>&1; then
        log_info "使用Trivy进行安全扫描..."
        trivy image "$image_name"
    elif command -v docker >/dev/null 2>&1 && docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy:latest --version >/dev/null 2>&1; then
        log_info "使用Docker版Trivy进行安全扫描..."
        docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy:latest image "$image_name"
    else
        log_warning "未找到安全扫描工具，跳过安全扫描"
    fi
}

# 显示帮助信息
show_help() {
    echo "Greatest Works Docker构建脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help              显示帮助信息"
    echo "  -r, --registry URL      镜像仓库地址"
    echo "  -n, --namespace NAME    命名空间"
    echo "  -t, --tag TAG           镜像标签"
    echo "  --platform PLATFORM     目标平台 (linux/amd64,linux/arm64)"
    echo "  --multi-arch            构建多架构镜像"
    echo "  --dev                   构建开发镜像"
    echo "  --push                  构建后推送镜像"
    echo "  --no-cache              不使用构建缓存"
    echo "  --scan                  运行安全扫描"
    echo "  --clean                 清理构建缓存"
    echo ""
    echo "环境变量:"
    echo "  REGISTRY_USERNAME       镜像仓库用户名"
    echo "  REGISTRY_PASSWORD       镜像仓库密码"
    echo "  DOCKER_BUILDKIT         启用BuildKit (1/0)"
    echo ""
    echo "示例:"
    echo "  $0                      # 构建本地镜像"
    echo "  $0 --multi-arch --push  # 构建多架构镜像并推送"
    echo "  $0 --dev                # 构建开发镜像"
    echo "  $0 --tag v1.0.0 --push  # 构建指定版本并推送"
    echo "  $0 --clean              # 清理构建缓存"
}

# 主函数
main() {
    local registry="$DEFAULT_REGISTRY"
    local namespace="$DEFAULT_NAMESPACE"
    local tag="$DEFAULT_TAG"
    local platforms="linux/amd64"
    local multi_arch=false
    local dev_build=false
    local push=false
    local no_cache=false
    local scan=false
    local clean_only=false
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -r|--registry)
                registry="$2"
                shift 2
                ;;
            -n|--namespace)
                namespace="$2"
                shift 2
                ;;
            -t|--tag)
                tag="$2"
                shift 2
                ;;
            --platform)
                platforms="$2"
                shift 2
                ;;
            --multi-arch)
                multi_arch=true
                platforms="linux/amd64,linux/arm64"
                shift
                ;;
            --dev)
                dev_build=true
                shift
                ;;
            --push)
                push=true
                shift
                ;;
            --no-cache)
                no_cache=true
                shift
                ;;
            --scan)
                scan=true
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
    
    log_info "Greatest Works Docker构建脚本启动"
    
    # 检查Docker环境
    check_docker
    
    # 清理缓存
    if [ "$clean_only" = true ]; then
        clean_build_cache
        exit 0
    fi
    
    # 获取版本信息
    get_version_info
    
    # 启用BuildKit
    export DOCKER_BUILDKIT=1
    
    # 设置no-cache参数
    if [ "$no_cache" = true ]; then
        export DOCKER_BUILDKIT_NO_CACHE=1
    fi
    
    # 开发镜像构建
    if [ "$dev_build" = true ]; then
        build_dev_image
        
        if [ "$scan" = true ]; then
            run_security_scan "$PROJECT_NAME:dev"
        fi
        
        exit 0
    fi
    
    # 登录镜像仓库
    if [ "$push" = true ]; then
        login_registry "$registry"
    fi
    
    # 多架构构建
    if [ "$multi_arch" = true ]; then
        check_buildx
        build_multi_arch "$registry" "$namespace" "$tag" "$platforms" "$push"
    else
        # 单架构构建
        if [ "$push" = true ]; then
            check_buildx
            build_single_arch "$registry" "$namespace" "$tag" "$platforms" "$push"
        else
            # 本地构建
            local image_name="$registry/$namespace/$PROJECT_NAME:$tag"
            
            docker build \
                --build-arg VERSION="$VERSION" \
                --build-arg GIT_COMMIT="$GIT_COMMIT" \
                --build-arg GIT_BRANCH="$GIT_BRANCH" \
                --build-arg BUILD_TIME="$BUILD_TIME" \
                -t "$image_name" \
                .
            
            if [ $? -eq 0 ]; then
                log_success "镜像构建成功: $image_name"
                
                # 显示镜像信息
                show_image_info "$image_name"
                
                # 安全扫描
                if [ "$scan" = true ]; then
                    run_security_scan "$image_name"
                fi
            else
                log_error "镜像构建失败"
                exit 1
            fi
        fi
    fi
    
    log_success "Docker构建完成！"
    
    # 显示后续操作提示
    if [ "$push" = false ]; then
        local image_name="$registry/$namespace/$PROJECT_NAME:$tag"
        log_info "后续操作:"
        log_info "  运行镜像: docker run -p 9080:9080 $image_name"
        log_info "  推送镜像: docker push $image_name"
        log_info "  安全扫描: $0 --scan"
    fi
}

# 执行主函数
main "$@"