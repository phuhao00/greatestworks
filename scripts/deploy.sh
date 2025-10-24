#!/bin/bash

# Greatest Works - 部署脚本
# 支持多环境自动化部署

set -e

# 默认配置
DEFAULT_ENV="development"
DEFAULT_REGISTRY="registry.example.com"
DEFAULT_NAMESPACE="gaming"
PROJECT_NAME="greatestworks"

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

# 检查必要的工具
check_dependencies() {
    local missing_tools=()
    
    # 检查Docker
    if ! command -v docker &> /dev/null; then
        missing_tools+=("docker")
    fi
    
    # 根据部署类型检查其他工具
    case $DEPLOY_TYPE in
        "kubernetes")
            if ! command -v kubectl &> /dev/null; then
                missing_tools+=("kubectl")
            fi
            ;;
        "docker-compose")
            if ! command -v docker-compose &> /dev/null; then
                missing_tools+=("docker-compose")
            fi
            ;;
    esac
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        log_error "缺少必要工具: ${missing_tools[*]}"
        exit 1
    fi
}

# 构建Docker镜像
build_docker_image() {
    local tag="$1"
    
    log_info "构建Docker镜像: $tag"
    
    # 检查Dockerfile是否存在
    if [ ! -f "Dockerfile" ]; then
        log_error "Dockerfile不存在"
        exit 1
    fi
    
    # 构建镜像
    docker build -t "$tag" .
    
    if [ $? -eq 0 ]; then
        log_success "Docker镜像构建完成: $tag"
    else
        log_error "Docker镜像构建失败"
        exit 1
    fi
}

# 推送Docker镜像
push_docker_image() {
    local tag="$1"
    
    log_info "推送Docker镜像: $tag"
    
    # 登录到镜像仓库
    if [ -n "$REGISTRY_USERNAME" ] && [ -n "$REGISTRY_PASSWORD" ]; then
        echo "$REGISTRY_PASSWORD" | docker login "$REGISTRY" -u "$REGISTRY_USERNAME" --password-stdin
    fi
    
    # 推送镜像
    docker push "$tag"
    
    if [ $? -eq 0 ]; then
        log_success "Docker镜像推送完成: $tag"
    else
        log_error "Docker镜像推送失败"
        exit 1
    fi
}

# Docker Compose部署
deploy_docker_compose() {
    local env="$1"
    local compose_file="docker-compose.yml"
    
    # 检查环境特定的compose文件
    if [ -f "docker-compose.${env}.yml" ]; then
        compose_file="docker-compose.${env}.yml"
    fi
    
    log_info "使用Docker Compose部署 (环境: $env)"
    log_info "Compose文件: $compose_file"
    
    # 检查compose文件是否存在
    if [ ! -f "$compose_file" ]; then
        log_error "Docker Compose文件不存在: $compose_file"
        exit 1
    fi
    
    # 设置环境变量
    export DEPLOY_ENV="$env"
    export IMAGE_TAG="${IMAGE_TAG:-latest}"
    
    # 停止现有服务
    log_info "停止现有服务..."
    docker-compose -f "$compose_file" down
    
    # 拉取最新镜像
    log_info "拉取最新镜像..."
    docker-compose -f "$compose_file" pull
    
    # 启动服务
    log_info "启动服务..."
    docker-compose -f "$compose_file" up -d
    
    # 等待服务启动
    sleep 10
    
    # 检查服务状态
    log_info "检查服务状态..."
    docker-compose -f "$compose_file" ps
    
    log_success "Docker Compose部署完成"
}

# Kubernetes部署
deploy_kubernetes() {
    local env="$1"
    local namespace="${NAMESPACE:-$DEFAULT_NAMESPACE}"
    local k8s_dir="k8s"
    
    log_info "使用Kubernetes部署 (环境: $env, 命名空间: $namespace)"
    
    # 检查k8s目录是否存在
    if [ ! -d "$k8s_dir" ]; then
        log_error "Kubernetes配置目录不存在: $k8s_dir"
        exit 1
    fi
    
    # 创建命名空间（如果不存在）
    kubectl create namespace "$namespace" --dry-run=client -o yaml | kubectl apply -f -
    
    # 设置当前上下文的命名空间
    kubectl config set-context --current --namespace="$namespace"
    
    # 应用配置文件
    local config_files=()
    
    # 查找环境特定的配置文件
    if [ -d "$k8s_dir/$env" ]; then
        config_files=("$k8s_dir/$env"/*.yaml "$k8s_dir/$env"/*.yml)
    else
        config_files=("$k8s_dir"/*.yaml "$k8s_dir"/*.yml)
    fi
    
    # 替换环境变量
    export DEPLOY_ENV="$env"
    export IMAGE_TAG="${IMAGE_TAG:-latest}"
    export NAMESPACE="$namespace"
    
    for config_file in "${config_files[@]}"; do
        if [ -f "$config_file" ]; then
            log_info "应用配置: $config_file"
            envsubst < "$config_file" | kubectl apply -f -
        fi
    done
    
    # 等待部署完成
    log_info "等待部署完成..."
    kubectl rollout status deployment/$PROJECT_NAME -n "$namespace" --timeout=300s
    
    # 显示部署状态
    log_info "部署状态:"
    kubectl get pods -n "$namespace" -l app=$PROJECT_NAME
    
    log_success "Kubernetes部署完成"
}

# 健康检查
health_check() {
    local url="$1"
    local max_attempts=30
    local attempt=1
    
    log_info "执行健康检查: $url"
    
    while [ $attempt -le $max_attempts ]; do
        if curl -f -s "$url" > /dev/null; then
            log_success "健康检查通过"
            return 0
        fi
        
        log_info "健康检查失败，重试 ($attempt/$max_attempts)..."
        sleep 10
        ((attempt++))
    done
    
    log_error "健康检查失败，服务可能未正常启动"
    return 1
}

# 回滚部署
rollback_deployment() {
    local env="$1"
    
    log_warning "开始回滚部署..."
    
    case $DEPLOY_TYPE in
        "kubernetes")
            local namespace="${NAMESPACE:-$DEFAULT_NAMESPACE}"
            kubectl rollout undo deployment/$PROJECT_NAME -n "$namespace"
            kubectl rollout status deployment/$PROJECT_NAME -n "$namespace"
            ;;
        "docker-compose")
            # Docker Compose回滚需要手动处理
            log_warning "Docker Compose回滚需要手动操作"
            ;;
    esac
    
    log_success "回滚完成"
}

# 显示帮助信息
show_help() {
    echo "Greatest Works 部署脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help              显示帮助信息"
    echo "  -e, --env ENV           部署环境 (development/staging/production)"
    echo "  -t, --type TYPE         部署类型 (docker-compose/kubernetes)"
    echo "  -r, --registry URL      Docker镜像仓库地址"
    echo "  -n, --namespace NAME    Kubernetes命名空间"
    echo "  --tag TAG               Docker镜像标签"
    echo "  --build                 构建Docker镜像"
    echo "  --push                  推送Docker镜像"
    echo "  --no-health-check       跳过健康检查"
    echo "  --rollback              回滚部署"
    echo ""
    echo "环境变量:"
    echo "  REGISTRY_USERNAME       镜像仓库用户名"
    echo "  REGISTRY_PASSWORD       镜像仓库密码"
    echo "  HEALTH_CHECK_URL        健康检查URL"
    echo ""
    echo "示例:"
    echo "  $0 -e development -t docker-compose"
    echo "  $0 -e production -t kubernetes --build --push"
    echo "  $0 --rollback -e production -t kubernetes"
}

# 主函数
main() {
    local env="$DEFAULT_ENV"
    local deploy_type="docker-compose"
    local registry="$DEFAULT_REGISTRY"
    local namespace="$DEFAULT_NAMESPACE"
    local image_tag="latest"
    local build_image=false
    local push_image=false
    local skip_health_check=false
    local rollback=false
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -e|--env)
                env="$2"
                shift 2
                ;;
            -t|--type)
                deploy_type="$2"
                shift 2
                ;;
            -r|--registry)
                registry="$2"
                shift 2
                ;;
            -n|--namespace)
                namespace="$2"
                shift 2
                ;;
            --tag)
                image_tag="$2"
                shift 2
                ;;
            --build)
                build_image=true
                shift
                ;;
            --push)
                push_image=true
                shift
                ;;
            --no-health-check)
                skip_health_check=true
                shift
                ;;
            --rollback)
                rollback=true
                shift
                ;;
            *)
                log_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 设置全局变量
    export DEPLOY_TYPE="$deploy_type"
    export REGISTRY="$registry"
    export NAMESPACE="$namespace"
    export IMAGE_TAG="$image_tag"
    
    log_info "Greatest Works 部署脚本启动"
    log_info "环境: $env"
    log_info "部署类型: $deploy_type"
    log_info "镜像标签: $image_tag"
    
    # 检查依赖
    check_dependencies
    
    # 处理回滚
    if [ "$rollback" = true ]; then
        rollback_deployment "$env"
        exit 0
    fi
    
    # 构建镜像
    if [ "$build_image" = true ]; then
        local full_tag="$registry/$PROJECT_NAME:$image_tag"
        build_docker_image "$full_tag"
        
        # 推送镜像
        if [ "$push_image" = true ]; then
            push_docker_image "$full_tag"
        fi
    fi
    
    # 执行部署
    case $deploy_type in
        "docker-compose")
            deploy_docker_compose "$env"
            ;;
        "kubernetes")
            deploy_kubernetes "$env"
            ;;
        *)
            log_error "不支持的部署类型: $deploy_type"
            exit 1
            ;;
    esac
    
    # 健康检查
    if [ "$skip_health_check" != true ]; then
        local health_url="${HEALTH_CHECK_URL:-http://localhost:8080/health}"
        if ! health_check "$health_url"; then
            log_error "部署可能失败，请检查服务状态"
            exit 1
        fi
    fi
    
    log_success "部署完成！"
}

# 执行主函数
main "$@"