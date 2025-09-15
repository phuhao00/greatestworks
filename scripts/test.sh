#!/bin/bash

# Greatest Works - 测试脚本
# 运行单元测试、集成测试和性能测试

set -e

# 默认配置
DEFAULT_TIMEOUT="30s"
DEFAULT_COVERAGE_THRESHOLD=80
TEST_RESULTS_DIR="./test-results"
COVERAGE_DIR="./coverage"

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

# 检查Go环境
check_go_env() {
    if ! command -v go &> /dev/null; then
        log_error "Go未安装或不在PATH中"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    log_info "Go版本: ${GO_VERSION}"
}

# 准备测试环境
setup_test_env() {
    log_info "准备测试环境..."
    
    # 创建测试结果目录
    mkdir -p "$TEST_RESULTS_DIR"
    mkdir -p "$COVERAGE_DIR"
    
    # 清理之前的测试结果
    rm -f "$TEST_RESULTS_DIR"/*.xml
    rm -f "$TEST_RESULTS_DIR"/*.json
    rm -f "$COVERAGE_DIR"/*.out
    rm -f "$COVERAGE_DIR"/*.html
    
    # 下载测试依赖
    go mod download
    
    log_success "测试环境准备完成"
}

# 运行单元测试
run_unit_tests() {
    local timeout="${TIMEOUT:-$DEFAULT_TIMEOUT}"
    local verbose="${VERBOSE:-false}"
    local race="${RACE:-true}"
    
    log_info "运行单元测试..."
    
    local test_args=()
    test_args+=("-timeout=$timeout")
    
    if [ "$verbose" = true ]; then
        test_args+=("-v")
    fi
    
    if [ "$race" = true ]; then
        test_args+=("-race")
    fi
    
    # 运行测试并生成覆盖率报告
    local coverage_file="$COVERAGE_DIR/unit.out"
    test_args+=("-coverprofile=$coverage_file")
    test_args+=("-covermode=atomic")
    
    # 排除某些包
    local packages=$(go list ./... | grep -v -E '(cmd/|scripts/|docs/|test/)')
    
    log_info "测试包: $(echo $packages | wc -w) 个"
    
    if go test "${test_args[@]}" $packages; then
        log_success "单元测试通过"
        
        # 生成覆盖率报告
        if [ -f "$coverage_file" ]; then
            generate_coverage_report "$coverage_file" "unit"
        fi
        
        return 0
    else
        log_error "单元测试失败"
        return 1
    fi
}

# 运行集成测试
run_integration_tests() {
    local timeout="${INTEGRATION_TIMEOUT:-60s}"
    
    log_info "运行集成测试..."
    
    # 检查是否有集成测试
    if ! find . -name "*_integration_test.go" -o -name "*_test.go" | grep -q integration; then
        log_warning "未找到集成测试文件"
        return 0
    fi
    
    # 启动测试依赖服务（如果需要）
    if [ -f "docker-compose.test.yml" ]; then
        log_info "启动测试依赖服务..."
        docker-compose -f docker-compose.test.yml up -d
        
        # 等待服务启动
        sleep 10
        
        # 设置清理函数
        trap 'docker-compose -f docker-compose.test.yml down' EXIT
    fi
    
    # 运行集成测试
    local coverage_file="$COVERAGE_DIR/integration.out"
    
    if go test -tags=integration -timeout="$timeout" -coverprofile="$coverage_file" ./...; then
        log_success "集成测试通过"
        
        # 生成覆盖率报告
        if [ -f "$coverage_file" ]; then
            generate_coverage_report "$coverage_file" "integration"
        fi
        
        return 0
    else
        log_error "集成测试失败"
        return 1
    fi
}

# 运行性能测试
run_benchmark_tests() {
    local benchtime="${BENCHTIME:-10s}"
    local benchmem="${BENCHMEM:-true}"
    
    log_info "运行性能测试..."
    
    # 检查是否有性能测试
    if ! find . -name "*_test.go" -exec grep -l "func Benchmark" {} \; | head -1 > /dev/null; then
        log_warning "未找到性能测试函数"
        return 0
    fi
    
    local bench_args=()
    bench_args+=("-bench=.")
    bench_args+=("-benchtime=$benchtime")
    
    if [ "$benchmem" = true ]; then
        bench_args+=("-benchmem")
    fi
    
    # 运行性能测试
    local bench_output="$TEST_RESULTS_DIR/benchmark.txt"
    
    if go test "${bench_args[@]}" ./... | tee "$bench_output"; then
        log_success "性能测试完成"
        log_info "性能测试结果保存到: $bench_output"
        return 0
    else
        log_error "性能测试失败"
        return 1
    fi
}

# 生成覆盖率报告
generate_coverage_report() {
    local coverage_file="$1"
    local test_type="$2"
    
    if [ ! -f "$coverage_file" ]; then
        log_warning "覆盖率文件不存在: $coverage_file"
        return
    fi
    
    log_info "生成${test_type}覆盖率报告..."
    
    # 生成HTML报告
    local html_file="$COVERAGE_DIR/${test_type}.html"
    go tool cover -html="$coverage_file" -o "$html_file"
    
    # 计算覆盖率百分比
    local coverage_percent=$(go tool cover -func="$coverage_file" | grep total | awk '{print $3}' | sed 's/%//')
    
    log_info "${test_type}测试覆盖率: ${coverage_percent}%"
    log_info "HTML报告: $html_file"
    
    # 检查覆盖率阈值
    local threshold="${COVERAGE_THRESHOLD:-$DEFAULT_COVERAGE_THRESHOLD}"
    if (( $(echo "$coverage_percent < $threshold" | bc -l) )); then
        log_warning "${test_type}测试覆盖率 (${coverage_percent}%) 低于阈值 (${threshold}%)"
        if [ "$FAIL_ON_LOW_COVERAGE" = true ]; then
            return 1
        fi
    else
        log_success "${test_type}测试覆盖率达标"
    fi
}

# 合并覆盖率报告
merge_coverage_reports() {
    log_info "合并覆盖率报告..."
    
    local coverage_files=()
    for file in "$COVERAGE_DIR"/*.out; do
        if [ -f "$file" ]; then
            coverage_files+=("$file")
        fi
    done
    
    if [ ${#coverage_files[@]} -eq 0 ]; then
        log_warning "未找到覆盖率文件"
        return
    fi
    
    # 合并覆盖率文件
    local merged_file="$COVERAGE_DIR/merged.out"
    
    # 写入mode行
    echo "mode: atomic" > "$merged_file"
    
    # 合并所有覆盖率数据
    for file in "${coverage_files[@]}"; do
        tail -n +2 "$file" >> "$merged_file"
    done
    
    # 生成合并后的HTML报告
    generate_coverage_report "$merged_file" "merged"
    
    log_success "覆盖率报告合并完成"
}

# 运行代码质量检查
run_quality_checks() {
    log_info "运行代码质量检查..."
    
    # golangci-lint检查
    if command -v golangci-lint &> /dev/null; then
        log_info "运行golangci-lint..."
        if golangci-lint run --out-format=json > "$TEST_RESULTS_DIR/lint.json"; then
            log_success "代码质量检查通过"
        else
            log_error "代码质量检查失败"
            return 1
        fi
    else
        log_warning "golangci-lint未安装，跳过代码质量检查"
    fi
    
    # go vet检查
    log_info "运行go vet..."
    if go vet ./...; then
        log_success "go vet检查通过"
    else
        log_error "go vet检查失败"
        return 1
    fi
    
    # go fmt检查
    log_info "检查代码格式..."
    local unformatted=$(gofmt -l .)
    if [ -n "$unformatted" ]; then
        log_error "以下文件格式不正确:"
        echo "$unformatted"
        return 1
    else
        log_success "代码格式检查通过"
    fi
}

# 生成测试报告
generate_test_report() {
    log_info "生成测试报告..."
    
    local report_file="$TEST_RESULTS_DIR/test-report.html"
    
    cat > "$report_file" << EOF
<!DOCTYPE html>
<html>
<head>
    <title>Greatest Works - 测试报告</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background-color: #f0f0f0; padding: 20px; border-radius: 5px; }
        .section { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 5px; }
        .success { background-color: #d4edda; border-color: #c3e6cb; }
        .warning { background-color: #fff3cd; border-color: #ffeaa7; }
        .error { background-color: #f8d7da; border-color: #f5c6cb; }
        .coverage { font-size: 18px; font-weight: bold; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Greatest Works - 测试报告</h1>
        <p>生成时间: $(date)</p>
        <p>Go版本: $(go version)</p>
    </div>
EOF
    
    # 添加覆盖率信息
    if [ -f "$COVERAGE_DIR/merged.out" ]; then
        local total_coverage=$(go tool cover -func="$COVERAGE_DIR/merged.out" | grep total | awk '{print $3}')
        echo "    <div class='section success'>" >> "$report_file"
        echo "        <h2>总体覆盖率</h2>" >> "$report_file"
        echo "        <p class='coverage'>$total_coverage</p>" >> "$report_file"
        echo "    </div>" >> "$report_file"
    fi
    
    echo "</body></html>" >> "$report_file"
    
    log_success "测试报告生成完成: $report_file"
}

# 清理测试环境
cleanup_test_env() {
    log_info "清理测试环境..."
    
    # 停止测试服务
    if [ -f "docker-compose.test.yml" ]; then
        docker-compose -f docker-compose.test.yml down
    fi
    
    log_success "测试环境清理完成"
}

# 显示帮助信息
show_help() {
    echo "Greatest Works 测试脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help              显示帮助信息"
    echo "  -u, --unit              只运行单元测试"
    echo "  -i, --integration       只运行集成测试"
    echo "  -b, --benchmark         只运行性能测试"
    echo "  -q, --quality           只运行代码质量检查"
    echo "  -a, --all               运行所有测试 (默认)"
    echo "  -v, --verbose           详细输出"
    echo "  --no-race               禁用竞态检测"
    echo "  --timeout DURATION      测试超时时间"
    echo "  --coverage-threshold N  覆盖率阈值"
    echo "  --fail-on-low-coverage  覆盖率不达标时失败"
    echo "  --clean                 清理测试结果"
    echo ""
    echo "环境变量:"
    echo "  VERBOSE                 详细输出 (true/false)"
    echo "  RACE                    竞态检测 (true/false)"
    echo "  TIMEOUT                 单元测试超时时间"
    echo "  INTEGRATION_TIMEOUT     集成测试超时时间"
    echo "  BENCHTIME               性能测试运行时间"
    echo "  COVERAGE_THRESHOLD      覆盖率阈值"
    echo ""
    echo "示例:"
    echo "  $0                      # 运行所有测试"
    echo "  $0 -u -v               # 运行单元测试，详细输出"
    echo "  $0 --coverage-threshold 90  # 设置覆盖率阈值为90%"
}

# 主函数
main() {
    local run_unit=false
    local run_integration=false
    local run_benchmark=false
    local run_quality=false
    local run_all=true
    local clean_only=false
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -u|--unit)
                run_unit=true
                run_all=false
                shift
                ;;
            -i|--integration)
                run_integration=true
                run_all=false
                shift
                ;;
            -b|--benchmark)
                run_benchmark=true
                run_all=false
                shift
                ;;
            -q|--quality)
                run_quality=true
                run_all=false
                shift
                ;;
            -a|--all)
                run_all=true
                shift
                ;;
            -v|--verbose)
                export VERBOSE=true
                shift
                ;;
            --no-race)
                export RACE=false
                shift
                ;;
            --timeout)
                export TIMEOUT="$2"
                shift 2
                ;;
            --coverage-threshold)
                export COVERAGE_THRESHOLD="$2"
                shift 2
                ;;
            --fail-on-low-coverage)
                export FAIL_ON_LOW_COVERAGE=true
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
    
    log_info "Greatest Works 测试脚本启动"
    
    # 检查Go环境
    check_go_env
    
    # 准备测试环境
    setup_test_env
    
    if [ "$clean_only" = true ]; then
        cleanup_test_env
        log_success "清理完成"
        exit 0
    fi
    
    # 设置清理函数
    trap cleanup_test_env EXIT
    
    local test_failed=false
    
    # 运行测试
    if [ "$run_all" = true ] || [ "$run_quality" = true ]; then
        if ! run_quality_checks; then
            test_failed=true
        fi
    fi
    
    if [ "$run_all" = true ] || [ "$run_unit" = true ]; then
        if ! run_unit_tests; then
            test_failed=true
        fi
    fi
    
    if [ "$run_all" = true ] || [ "$run_integration" = true ]; then
        if ! run_integration_tests; then
            test_failed=true
        fi
    fi
    
    if [ "$run_all" = true ] || [ "$run_benchmark" = true ]; then
        if ! run_benchmark_tests; then
            test_failed=true
        fi
    fi
    
    # 合并覆盖率报告
    if [ "$run_all" = true ] || [ "$run_unit" = true ] || [ "$run_integration" = true ]; then
        merge_coverage_reports
    fi
    
    # 生成测试报告
    generate_test_report
    
    if [ "$test_failed" = true ]; then
        log_error "部分测试失败"
        exit 1
    else
        log_success "所有测试通过！"
    fi
}

# 执行主函数
main "$@"