#!/bin/bash

# Greatest Works - 开发环境设置脚本
# 一键设置完整的开发环境

set -e

# 默认配置
GO_VERSION="1.21"
NODE_VERSION="18"
DOCKER_COMPOSE_VERSION="2.20.0"

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

# 检测操作系统
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        OS="linux"
        if [ -f /etc/debian_version ]; then
            DISTRO="debian"
        elif [ -f /etc/redhat-release ]; then
            DISTRO="redhat"
        else
            DISTRO="unknown"
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macos"
        DISTRO="macos"
    elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]]; then
        OS="windows"
        DISTRO="windows"
    else
        OS="unknown"
        DISTRO="unknown"
    fi
    
    log_info "检测到操作系统: $OS ($DISTRO)"
}

# 检查命令是否存在
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# 检查Go环境
check_go() {
    log_info "检查Go环境..."
    
    if command_exists go; then
        local current_version=$(go version | awk '{print $3}' | sed 's/go//')
        log_info "当前Go版本: $current_version"
        
        # 检查版本是否满足要求
        if printf '%s\n%s\n' "$GO_VERSION" "$current_version" | sort -V -C; then
            log_success "Go版本满足要求"
            return 0
        else
            log_warning "Go版本过低，建议升级到 $GO_VERSION 或更高版本"
        fi
    else
        log_warning "Go未安装"
        return 1
    fi
}

# 安装Go
install_go() {
    log_info "安装Go $GO_VERSION..."
    
    case $OS in
        "linux")
            local arch=$(uname -m)
            case $arch in
                "x86_64") arch="amd64" ;;
                "aarch64") arch="arm64" ;;
                *) log_error "不支持的架构: $arch"; return 1 ;;
            esac
            
            local go_url="https://golang.org/dl/go${GO_VERSION}.linux-${arch}.tar.gz"
            local go_file="/tmp/go${GO_VERSION}.linux-${arch}.tar.gz"
            
            curl -L "$go_url" -o "$go_file"
            sudo rm -rf /usr/local/go
            sudo tar -C /usr/local -xzf "$go_file"
            rm "$go_file"
            
            # 添加到PATH
            if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
                echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
            fi
            ;;
        "macos")
            if command_exists brew; then
                brew install go
            else
                log_error "请先安装Homebrew或手动安装Go"
                return 1
            fi
            ;;
        "windows")
            log_error "Windows环境请手动安装Go"
            return 1
            ;;
        *)
            log_error "不支持的操作系统: $OS"
            return 1
            ;;
    esac
    
    log_success "Go安装完成"
}

# 检查Docker环境
check_docker() {
    log_info "检查Docker环境..."
    
    if command_exists docker; then
        local docker_version=$(docker --version | awk '{print $3}' | sed 's/,//')
        log_info "当前Docker版本: $docker_version"
        
        # 检查Docker是否运行
        if docker info >/dev/null 2>&1; then
            log_success "Docker运行正常"
        else
            log_warning "Docker未运行，请启动Docker服务"
        fi
    else
        log_warning "Docker未安装"
        return 1
    fi
    
    # 检查Docker Compose
    if command_exists docker-compose; then
        local compose_version=$(docker-compose --version | awk '{print $3}' | sed 's/,//')
        log_info "当前Docker Compose版本: $compose_version"
        log_success "Docker Compose已安装"
    else
        log_warning "Docker Compose未安装"
        return 1
    fi
}

# 安装Docker
install_docker() {
    log_info "安装Docker..."
    
    case $OS in
        "linux")
            case $DISTRO in
                "debian")
                    # 安装依赖
                    sudo apt-get update
                    sudo apt-get install -y apt-transport-https ca-certificates curl gnupg lsb-release
                    
                    # 添加Docker GPG密钥
                    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
                    
                    # 添加Docker仓库
                    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
                    
                    # 安装Docker
                    sudo apt-get update
                    sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
                    
                    # 添加用户到docker组
                    sudo usermod -aG docker $USER
                    ;;
                "redhat")
                    sudo yum install -y yum-utils
                    sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
                    sudo yum install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
                    sudo systemctl start docker
                    sudo systemctl enable docker
                    sudo usermod -aG docker $USER
                    ;;
                *)
                    log_error "不支持的Linux发行版: $DISTRO"
                    return 1
                    ;;
            esac
            ;;
        "macos")
            if command_exists brew; then
                brew install --cask docker
                log_info "请手动启动Docker Desktop"
            else
                log_error "请先安装Homebrew或手动安装Docker Desktop"
                return 1
            fi
            ;;
        "windows")
            log_error "Windows环境请手动安装Docker Desktop"
            return 1
            ;;
        *)
            log_error "不支持的操作系统: $OS"
            return 1
            ;;
    esac
    
    log_success "Docker安装完成"
    log_warning "请重新登录以使用户组更改生效"
}

# 检查开发工具
check_dev_tools() {
    log_info "检查开发工具..."
    
    local missing_tools=()
    
    # 检查git
    if ! command_exists git; then
        missing_tools+=("git")
    fi
    
    # 检查make
    if ! command_exists make; then
        missing_tools+=("make")
    fi
    
    # 检查curl
    if ! command_exists curl; then
        missing_tools+=("curl")
    fi
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        log_warning "缺少开发工具: ${missing_tools[*]}"
        return 1
    else
        log_success "开发工具检查通过"
        return 0
    fi
}

# 安装开发工具
install_dev_tools() {
    log_info "安装开发工具..."
    
    case $OS in
        "linux")
            case $DISTRO in
                "debian")
                    sudo apt-get update
                    sudo apt-get install -y git make curl wget unzip build-essential
                    ;;
                "redhat")
                    sudo yum groupinstall -y "Development Tools"
                    sudo yum install -y git make curl wget unzip
                    ;;
                *)
                    log_error "不支持的Linux发行版: $DISTRO"
                    return 1
                    ;;
            esac
            ;;
        "macos")
            if command_exists brew; then
                brew install git make curl wget
            else
                # 安装Xcode命令行工具
                xcode-select --install
            fi
            ;;
        "windows")
            log_error "Windows环境请手动安装开发工具"
            return 1
            ;;
        *)
            log_error "不支持的操作系统: $OS"
            return 1
            ;;
    esac
    
    log_success "开发工具安装完成"
}

# 安装Go开发工具
install_go_tools() {
    log_info "安装Go开发工具..."
    
    # 检查Go是否可用
    if ! command_exists go; then
        log_error "Go未安装，请先安装Go"
        return 1
    fi
    
    # 安装常用Go工具
    local go_tools=(
        "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
        "github.com/air-verse/air@latest"
        "github.com/swaggo/swag/cmd/swag@latest"
        "github.com/golang/mock/mockgen@latest"
        "golang.org/x/tools/cmd/goimports@latest"
        "golang.org/x/tools/gopls@latest"
    )
    
    for tool in "${go_tools[@]}"; do
        log_info "安装 $tool..."
        if go install "$tool"; then
            log_success "$tool 安装成功"
        else
            log_warning "$tool 安装失败"
        fi
    done
    
    # 确保GOPATH/bin在PATH中
    local gopath=$(go env GOPATH)
    if [ -n "$gopath" ] && [[ ":$PATH:" != *":$gopath/bin:"* ]]; then
        echo "export PATH=\$PATH:$gopath/bin" >> ~/.bashrc
        log_info "已添加 $gopath/bin 到PATH"
    fi
    
    log_success "Go开发工具安装完成"
}

# 设置项目环境
setup_project() {
    log_info "设置项目环境..."
    
    # 检查是否在项目根目录
    if [ ! -f "go.mod" ]; then
        log_error "请在项目根目录运行此脚本"
        return 1
    fi
    
    # 下载Go依赖
    log_info "下载Go模块依赖..."
    go mod download
    go mod tidy
    
    # 创建必要的目录
    local dirs=("bin" "logs" "tmp" "test-results" "coverage")
    for dir in "${dirs[@]}"; do
        if [ ! -d "$dir" ]; then
            mkdir -p "$dir"
            log_info "创建目录: $dir"
        fi
    done
    
    # 复制配置文件模板
    if [ ! -f "config.yaml" ]; then
        if [ -f "configs/config.dev.yaml.example" ]; then
            cp "configs/config.dev.yaml.example" "config.yaml"
            log_info "创建开发配置文件: config.yaml"
        fi
    fi
    
    # 创建.env文件
    if [ ! -f ".env" ]; then
        cat > .env << EOF
# Greatest Works 开发环境变量

# 服务器配置
SERVER_PORT=8080
SERVER_HOST=localhost

# 数据库配置
MONGODB_URI=mongodb://localhost:27017/gamedb_dev
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

# 消息队列
NATS_URL=nats://localhost:4222

# 认证配置
JWT_SECRET=dev-secret-key-change-in-production

# 日志配置
LOG_LEVEL=debug
LOG_FORMAT=text

# 开发模式
DEV_MODE=true
HOT_RELOAD=true
EOF
        log_info "创建环境变量文件: .env"
    fi
    
    # 设置Git hooks
    if [ -d ".git" ]; then
        setup_git_hooks
    fi
    
    log_success "项目环境设置完成"
}

# 设置Git hooks
setup_git_hooks() {
    log_info "设置Git hooks..."
    
    local hooks_dir=".git/hooks"
    
    # pre-commit hook
    cat > "$hooks_dir/pre-commit" << 'EOF'
#!/bin/bash

# Greatest Works pre-commit hook

set -e

echo "运行pre-commit检查..."

# 格式化代码
echo "格式化Go代码..."
gofmt -w .

# 运行代码质量检查
if command -v golangci-lint >/dev/null 2>&1; then
    echo "运行golangci-lint..."
    golangci-lint run
fi

# 运行测试
echo "运行单元测试..."
go test -short ./...

echo "pre-commit检查通过"
EOF
    
    chmod +x "$hooks_dir/pre-commit"
    log_info "设置pre-commit hook"
    
    log_success "Git hooks设置完成"
}

# 启动开发服务
start_dev_services() {
    log_info "启动开发服务..."
    
    # 检查docker-compose文件
    local compose_file="docker-compose.dev.yml"
    if [ ! -f "$compose_file" ]; then
        compose_file="docker-compose.yml"
    fi
    
    if [ -f "$compose_file" ]; then
        log_info "启动Docker服务..."
        docker-compose -f "$compose_file" up -d
        
        # 等待服务启动
        log_info "等待服务启动..."
        sleep 10
        
        # 检查服务状态
        docker-compose -f "$compose_file" ps
        
        log_success "开发服务启动完成"
    else
        log_warning "未找到docker-compose文件，跳过服务启动"
    fi
}

# 显示帮助信息
show_help() {
    echo "Greatest Works 开发环境设置脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help          显示帮助信息"
    echo "  --check-only        只检查环境，不安装"
    echo "  --skip-docker       跳过Docker安装"
    echo "  --skip-go-tools     跳过Go工具安装"
    echo "  --skip-services     跳过开发服务启动"
    echo "  --force-install     强制重新安装所有组件"
    echo ""
    echo "功能:"
    echo "  - 检查和安装Go环境"
    echo "  - 检查和安装Docker环境"
    echo "  - 安装开发工具和Go工具"
    echo "  - 设置项目环境"
    echo "  - 启动开发服务"
    echo ""
    echo "示例:"
    echo "  $0                  # 完整设置开发环境"
    echo "  $0 --check-only     # 只检查环境状态"
    echo "  $0 --skip-docker    # 跳过Docker相关设置"
}

# 主函数
main() {
    local check_only=false
    local skip_docker=false
    local skip_go_tools=false
    local skip_services=false
    local force_install=false
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            --check-only)
                check_only=true
                shift
                ;;
            --skip-docker)
                skip_docker=true
                shift
                ;;
            --skip-go-tools)
                skip_go_tools=true
                shift
                ;;
            --skip-services)
                skip_services=true
                shift
                ;;
            --force-install)
                force_install=true
                shift
                ;;
            *)
                log_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    log_info "Greatest Works 开发环境设置脚本启动"
    
    # 检测操作系统
    detect_os
    
    # 检查基础开发工具
    if ! check_dev_tools; then
        if [ "$check_only" = false ]; then
            install_dev_tools
        fi
    fi
    
    # 检查Go环境
    if ! check_go || [ "$force_install" = true ]; then
        if [ "$check_only" = false ]; then
            install_go
        fi
    fi
    
    # 检查Docker环境
    if [ "$skip_docker" = false ]; then
        if ! check_docker || [ "$force_install" = true ]; then
            if [ "$check_only" = false ]; then
                install_docker
            fi
        fi
    fi
    
    if [ "$check_only" = true ]; then
        log_info "环境检查完成"
        exit 0
    fi
    
    # 安装Go开发工具
    if [ "$skip_go_tools" = false ]; then
        install_go_tools
    fi
    
    # 设置项目环境
    setup_project
    
    # 启动开发服务
    if [ "$skip_services" = false ]; then
        start_dev_services
    fi
    
    log_success "开发环境设置完成！"
    log_info "请运行以下命令使环境变量生效:"
    log_info "  source ~/.bashrc"
    log_info ""
    log_info "然后可以使用以下命令:"
    log_info "  make dev          # 启动开发服务器"
    log_info "  make test         # 运行测试"
    log_info "  make build        # 构建项目"
}

# 执行主函数
main "$@"