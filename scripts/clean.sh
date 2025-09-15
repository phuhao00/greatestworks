#!/bin/bash

# Greatest Works - 清理脚本
# 清理构建产物、临时文件、日志文件等

set -e

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

# 清理构建产物
clean_build_artifacts() {
    log_info "清理构建产物..."
    
    local cleaned_items=()
    
    # 清理bin目录
    if [ -d "bin" ]; then
        rm -rf bin/*
        cleaned_items+=("bin目录")
    fi
    
    # 清理可执行文件
    local executables=("server" "greatestworks" "greatestworks.exe")
    for exe in "${executables[@]}"; do
        if [ -f "$exe" ]; then
            rm -f "$exe"
            cleaned_items+=("可执行文件: $exe")
        fi
    done
    
    # 清理Go构建缓存
    if command -v go >/dev/null 2>&1; then
        go clean -cache
        go clean -modcache
        cleaned_items+=("Go构建缓存")
    fi
    
    if [ ${#cleaned_items[@]} -gt 0 ]; then
        log_success "构建产物清理完成: ${cleaned_items[*]}"
    else
        log_info "没有找到构建产物"
    fi
}

# 清理临时文件
clean_temp_files() {
    log_info "清理临时文件..."
    
    local cleaned_items=()
    local temp_patterns=(
        "*.tmp"
        "*.temp"
        "*.swp"
        "*.swo"
        "*~"
        ".DS_Store"
        "Thumbs.db"
        "*.pid"
        "*.lock"
    )
    
    for pattern in "${temp_patterns[@]}"; do
        local files=$(find . -name "$pattern" -type f 2>/dev/null || true)
        if [ -n "$files" ]; then
            echo "$files" | xargs rm -f
            cleaned_items+=("$pattern")
        fi
    done
    
    # 清理tmp目录
    if [ -d "tmp" ]; then
        rm -rf tmp/*
        cleaned_items+=("tmp目录")
    fi
    
    # 清理系统临时文件
    local temp_dirs=("/tmp/greatestworks_*" "/tmp/build_*" "/tmp/test_*")
    for temp_dir in "${temp_dirs[@]}"; do
        if ls $temp_dir 1> /dev/null 2>&1; then
            rm -rf $temp_dir
            cleaned_items+=("系统临时文件")
        fi
    done
    
    if [ ${#cleaned_items[@]} -gt 0 ]; then
        log_success "临时文件清理完成: ${cleaned_items[*]}"
    else
        log_info "没有找到临时文件"
    fi
}

# 清理日志文件
clean_log_files() {
    log_info "清理日志文件..."
    
    local cleaned_items=()
    
    # 清理logs目录
    if [ -d "logs" ]; then
        local log_count=$(find logs -name "*.log" -type f | wc -l)
        if [ $log_count -gt 0 ]; then
            find logs -name "*.log" -type f -delete
            cleaned_items+=("$log_count 个日志文件")
        fi
        
        # 清理压缩的日志文件
        local gz_count=$(find logs -name "*.log.gz" -type f | wc -l)
        if [ $gz_count -gt 0 ]; then
            find logs -name "*.log.gz" -type f -delete
            cleaned_items+=("$gz_count 个压缩日志文件")
        fi
    fi
    
    # 清理根目录的日志文件
    local root_logs=("*.log" "nohup.out" "error.log" "access.log")
    for log_pattern in "${root_logs[@]}"; do
        if ls $log_pattern 1> /dev/null 2>&1; then
            rm -f $log_pattern
            cleaned_items+=("根目录日志: $log_pattern")
        fi
    done
    
    if [ ${#cleaned_items[@]} -gt 0 ]; then
        log_success "日志文件清理完成: ${cleaned_items[*]}"
    else
        log_info "没有找到日志文件"
    fi
}

# 清理测试文件
clean_test_files() {
    log_info "清理测试文件..."
    
    local cleaned_items=()
    
    # 清理测试结果目录
    if [ -d "test-results" ]; then
        rm -rf test-results/*
        cleaned_items+=("test-results目录")
    fi
    
    # 清理覆盖率文件
    if [ -d "coverage" ]; then
        rm -rf coverage/*
        cleaned_items+=("coverage目录")
    fi
    
    # 清理Go测试缓存
    if command -v go >/dev/null 2>&1; then
        go clean -testcache
        cleaned_items+=("Go测试缓存")
    fi
    
    # 清理测试相关文件
    local test_files=("*.test" "*.out" "*.prof" "cpu.prof" "mem.prof")
    for test_file in "${test_files[@]}"; do
        if ls $test_file 1> /dev/null 2>&1; then
            rm -f $test_file
            cleaned_items+=("测试文件: $test_file")
        fi
    done
    
    if [ ${#cleaned_items[@]} -gt 0 ]; then
        log_success "测试文件清理完成: ${cleaned_items[*]}"
    else
        log_info "没有找到测试文件"
    fi
}

# 清理Docker资源
clean_docker_resources() {
    log_info "清理Docker资源..."
    
    if ! command -v docker >/dev/null 2>&1; then
        log_warning "Docker未安装，跳过Docker清理"
        return 0
    fi
    
    if ! docker info >/dev/null 2>&1; then
        log_warning "Docker未运行，跳过Docker清理"
        return 0
    fi
    
    local cleaned_items=()
    
    # 清理悬空镜像
    local dangling_images=$(docker images -f "dangling=true" -q)
    if [ -n "$dangling_images" ]; then
        echo "$dangling_images" | xargs docker rmi
        cleaned_items+=("悬空镜像")
    fi
    
    # 清理未使用的镜像
    docker image prune -f >/dev/null 2>&1
    cleaned_items+=("未使用镜像")
    
    # 清理停止的容器
    local stopped_containers=$(docker ps -a -q -f status=exited)
    if [ -n "$stopped_containers" ]; then
        echo "$stopped_containers" | xargs docker rm
        cleaned_items+=("停止的容器")
    fi
    
    # 清理未使用的网络
    docker network prune -f >/dev/null 2>&1
    cleaned_items+=("未使用网络")
    
    # 清理未使用的卷
    docker volume prune -f >/dev/null 2>&1
    cleaned_items+=("未使用卷")
    
    # 清理构建缓存
    if command -v docker >/dev/null 2>&1 && docker buildx version >/dev/null 2>&1; then
        docker buildx prune -f >/dev/null 2>&1
        cleaned_items+=("构建缓存")
    fi
    
    if [ ${#cleaned_items[@]} -gt 0 ]; then
        log_success "Docker资源清理完成: ${cleaned_items[*]}"
    else
        log_info "没有找到需要清理的Docker资源"
    fi
}

# 清理依赖缓存
clean_dependency_cache() {
    log_info "清理依赖缓存..."
    
    local cleaned_items=()
    
    # 清理Go模块缓存
    if command -v go >/dev/null 2>&1; then
        go clean -modcache
        cleaned_items+=("Go模块缓存")
    fi
    
    # 清理npm缓存
    if command -v npm >/dev/null 2>&1; then
        npm cache clean --force >/dev/null 2>&1
        cleaned_items+=("npm缓存")
    fi
    
    # 清理yarn缓存
    if command -v yarn >/dev/null 2>&1; then
        yarn cache clean >/dev/null 2>&1
        cleaned_items+=("yarn缓存")
    fi
    
    # 清理pip缓存
    if command -v pip >/dev/null 2>&1; then
        pip cache purge >/dev/null 2>&1
        cleaned_items+=("pip缓存")
    fi
    
    if [ ${#cleaned_items[@]} -gt 0 ]; then
        log_success "依赖缓存清理完成: ${cleaned_items[*]}"
    else
        log_info "没有找到依赖缓存"
    fi
}

# 清理IDE和编辑器文件
clean_ide_files() {
    log_info "清理IDE和编辑器文件..."
    
    local cleaned_items=()
    local ide_patterns=(
        ".vscode/settings.json"
        ".idea/workspace.xml"
        ".idea/tasks.xml"
        "*.sublime-workspace"
        ".project"
        ".classpath"
    )
    
    for pattern in "${ide_patterns[@]}"; do
        if ls $pattern 1> /dev/null 2>&1; then
            rm -f $pattern
            cleaned_items+=("IDE文件: $pattern")
        fi
    done
    
    # 清理编辑器备份文件
    find . -name "*~" -type f -delete 2>/dev/null || true
    find . -name "*.swp" -type f -delete 2>/dev/null || true
    find . -name "*.swo" -type f -delete 2>/dev/null || true
    
    if [ ${#cleaned_items[@]} -gt 0 ]; then
        log_success "IDE文件清理完成: ${cleaned_items[*]}"
    else
        log_info "没有找到IDE文件"
    fi
}

# 清理备份文件
clean_backup_files() {
    local keep_days="${1:-7}"
    
    log_info "清理 $keep_days 天前的备份文件..."
    
    local cleaned_items=()
    
    # 清理数据库备份
    if [ -d "backups" ]; then
        local old_backups=$(find backups -name "*.tar.gz" -type f -mtime +$keep_days 2>/dev/null || true)
        if [ -n "$old_backups" ]; then
            echo "$old_backups" | xargs rm -f
            local count=$(echo "$old_backups" | wc -l)
            cleaned_items+=("$count 个旧备份文件")
        fi
    fi
    
    # 清理日志备份
    if [ -d "logs" ]; then
        local old_log_backups=$(find logs -name "*.log.*.gz" -type f -mtime +$keep_days 2>/dev/null || true)
        if [ -n "$old_log_backups" ]; then
            echo "$old_log_backups" | xargs rm -f
            local count=$(echo "$old_log_backups" | wc -l)
            cleaned_items+=("$count 个旧日志备份")
        fi
    fi
    
    if [ ${#cleaned_items[@]} -gt 0 ]; then
        log_success "备份文件清理完成: ${cleaned_items[*]}"
    else
        log_info "没有找到需要清理的备份文件"
    fi
}

# 显示磁盘使用情况
show_disk_usage() {
    log_info "磁盘使用情况:"
    
    # 显示项目目录大小
    local project_size=$(du -sh . 2>/dev/null | cut -f1)
    echo "  项目总大小: $project_size"
    
    # 显示各个子目录大小
    local dirs=("bin" "logs" "tmp" "test-results" "coverage" "backups" "node_modules" ".git")
    for dir in "${dirs[@]}"; do
        if [ -d "$dir" ]; then
            local dir_size=$(du -sh "$dir" 2>/dev/null | cut -f1)
            echo "  $dir: $dir_size"
        fi
    done
    
    # 显示系统磁盘使用情况
    echo ""
    echo "系统磁盘使用情况:"
    df -h . | tail -1 | awk '{print "  可用空间: " $4 " / " $2 " (" $5 " 已使用)"}'
}

# 显示帮助信息
show_help() {
    echo "Greatest Works 清理脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help              显示帮助信息"
    echo "  -a, --all               清理所有内容 (默认)"
    echo "  -b, --build             只清理构建产物"
    echo "  -t, --temp              只清理临时文件"
    echo "  -l, --logs              只清理日志文件"
    echo "  --test                  只清理测试文件"
    echo "  --docker                只清理Docker资源"
    echo "  --cache                 只清理依赖缓存"
    echo "  --ide                   只清理IDE文件"
    echo "  --backup [DAYS]         清理N天前的备份文件 (默认7天)"
    echo "  --dry-run               显示将要清理的内容，但不执行"
    echo "  --usage                 显示磁盘使用情况"
    echo ""
    echo "示例:"
    echo "  $0                      # 清理所有内容"
    echo "  $0 --build              # 只清理构建产物"
    echo "  $0 --backup 30          # 清理30天前的备份"
    echo "  $0 --dry-run            # 预览清理内容"
    echo "  $0 --usage              # 显示磁盘使用情况"
}

# 主函数
main() {
    local clean_all=true
    local clean_build=false
    local clean_temp=false
    local clean_logs=false
    local clean_test=false
    local clean_docker=false
    local clean_cache=false
    local clean_ide=false
    local clean_backup=false
    local backup_days=7
    local dry_run=false
    local show_usage=false
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -a|--all)
                clean_all=true
                shift
                ;;
            -b|--build)
                clean_all=false
                clean_build=true
                shift
                ;;
            -t|--temp)
                clean_all=false
                clean_temp=true
                shift
                ;;
            -l|--logs)
                clean_all=false
                clean_logs=true
                shift
                ;;
            --test)
                clean_all=false
                clean_test=true
                shift
                ;;
            --docker)
                clean_all=false
                clean_docker=true
                shift
                ;;
            --cache)
                clean_all=false
                clean_cache=true
                shift
                ;;
            --ide)
                clean_all=false
                clean_ide=true
                shift
                ;;
            --backup)
                clean_all=false
                clean_backup=true
                if [[ $2 =~ ^[0-9]+$ ]]; then
                    backup_days=$2
                    shift 2
                else
                    shift
                fi
                ;;
            --dry-run)
                dry_run=true
                shift
                ;;
            --usage)
                show_usage=true
                shift
                ;;
            *)
                log_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    log_info "Greatest Works 清理脚本启动"
    
    # 显示磁盘使用情况
    if [ "$show_usage" = true ]; then
        show_disk_usage
        exit 0
    fi
    
    # 干运行模式
    if [ "$dry_run" = true ]; then
        log_warning "干运行模式：只显示将要清理的内容，不执行实际清理"
        # 在干运行模式下，可以添加检查逻辑
        exit 0
    fi
    
    # 执行清理操作
    if [ "$clean_all" = true ]; then
        clean_build_artifacts
        clean_temp_files
        clean_log_files
        clean_test_files
        clean_docker_resources
        clean_dependency_cache
        clean_ide_files
        clean_backup_files 7
    else
        if [ "$clean_build" = true ]; then
            clean_build_artifacts
        fi
        
        if [ "$clean_temp" = true ]; then
            clean_temp_files
        fi
        
        if [ "$clean_logs" = true ]; then
            clean_log_files
        fi
        
        if [ "$clean_test" = true ]; then
            clean_test_files
        fi
        
        if [ "$clean_docker" = true ]; then
            clean_docker_resources
        fi
        
        if [ "$clean_cache" = true ]; then
            clean_dependency_cache
        fi
        
        if [ "$clean_ide" = true ]; then
            clean_ide_files
        fi
        
        if [ "$clean_backup" = true ]; then
            clean_backup_files "$backup_days"
        fi
    fi
    
    # 显示清理后的磁盘使用情况
    echo ""
    show_disk_usage
    
    log_success "清理完成！"
}

# 执行主函数
main "$@"