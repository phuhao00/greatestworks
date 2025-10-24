#!/bin/bash

# Greatest Works - 数据库迁移脚本
# 管理MongoDB和Redis的数据迁移和初始化

set -e

# 默认配置
DEFAULT_MONGODB_URI="mongodb://localhost:27017"
DEFAULT_DATABASE="gamedb"
MIGRATIONS_DIR="./migrations"
SEEDS_DIR="./seeds"

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

# 检查MongoDB连接
check_mongodb() {
    local uri="${MONGODB_URI:-$DEFAULT_MONGODB_URI}"
    
    log_info "检查MongoDB连接: $uri"
    
    if command -v mongosh >/dev/null 2>&1; then
        if mongosh "$uri" --eval "db.adminCommand('ping')" >/dev/null 2>&1; then
            log_success "MongoDB连接正常"
            return 0
        else
            log_error "MongoDB连接失败"
            return 1
        fi
    elif command -v mongo >/dev/null 2>&1; then
        if mongo "$uri" --eval "db.adminCommand('ping')" >/dev/null 2>&1; then
            log_success "MongoDB连接正常"
            return 0
        else
            log_error "MongoDB连接失败"
            return 1
        fi
    else
        log_error "MongoDB客户端未安装 (mongosh 或 mongo)"
        return 1
    fi
}

# 检查Redis连接
check_redis() {
    local addr="${REDIS_ADDR:-localhost:6379}"
    
    log_info "检查Redis连接: $addr"
    
    if command -v redis-cli >/dev/null 2>&1; then
        if redis-cli -h "${addr%:*}" -p "${addr#*:}" ping >/dev/null 2>&1; then
            log_success "Redis连接正常"
            return 0
        else
            log_error "Redis连接失败"
            return 1
        fi
    else
        log_error "Redis客户端未安装 (redis-cli)"
        return 1
    fi
}

# 创建迁移目录
setup_migration_dirs() {
    log_info "创建迁移目录..."
    
    mkdir -p "$MIGRATIONS_DIR"
    mkdir -p "$SEEDS_DIR"
    
    # 创建迁移记录集合的初始化脚本
    if [ ! -f "$MIGRATIONS_DIR/000_init_migrations.js" ]; then
        cat > "$MIGRATIONS_DIR/000_init_migrations.js" << 'EOF'
// 初始化迁移记录集合
db.migrations.createIndex({ "version": 1 }, { unique: true });
db.migrations.createIndex({ "applied_at": 1 });

print("迁移记录集合初始化完成");
EOF
        log_info "创建初始化迁移文件"
    fi
    
    log_success "迁移目录设置完成"
}

# 生成新的迁移文件
generate_migration() {
    local name="$1"
    
    if [ -z "$name" ]; then
        log_error "请提供迁移文件名"
        return 1
    fi
    
    local timestamp=$(date +"%Y%m%d%H%M%S")
    local filename="${timestamp}_${name}.js"
    local filepath="$MIGRATIONS_DIR/$filename"
    
    cat > "$filepath" << EOF
// 迁移: $name
// 创建时间: $(date)
// 版本: $timestamp

// 执行迁移
function up() {
    print("执行迁移: $name");
    
    // 在这里添加迁移逻辑
    // 例如:
    // db.collection.createIndex({ "field": 1 });
    // db.collection.updateMany({}, { \$set: { "new_field": "default_value" } });
    
    print("迁移 $name 完成");
}

// 回滚迁移
function down() {
    print("回滚迁移: $name");
    
    // 在这里添加回滚逻辑
    // 例如:
    // db.collection.dropIndex({ "field": 1 });
    // db.collection.updateMany({}, { \$unset: { "new_field": "" } });
    
    print("回滚 $name 完成");
}

// 执行迁移
up();

// 记录迁移
db.migrations.insertOne({
    version: "$timestamp",
    name: "$name",
    filename: "$filename",
    applied_at: new Date()
});
EOF
    
    log_success "生成迁移文件: $filepath"
}

# 获取已应用的迁移
get_applied_migrations() {
    local uri="${MONGODB_URI:-$DEFAULT_MONGODB_URI}"
    local database="${DATABASE:-$DEFAULT_DATABASE}"
    
    if command -v mongosh >/dev/null 2>&1; then
        mongosh "$uri/$database" --quiet --eval "db.migrations.find({}, {version: 1, _id: 0}).sort({version: 1}).forEach(function(doc) { print(doc.version); })"
    else
        mongo "$uri/$database" --quiet --eval "db.migrations.find({}, {version: 1, _id: 0}).sort({version: 1}).forEach(function(doc) { print(doc.version); })"
    fi
}

# 获取待应用的迁移
get_pending_migrations() {
    local applied_migrations=()
    
    # 获取已应用的迁移
    while IFS= read -r version; do
        if [ -n "$version" ]; then
            applied_migrations+=("$version")
        fi
    done < <(get_applied_migrations)
    
    # 获取所有迁移文件
    local all_migrations=()
    for file in "$MIGRATIONS_DIR"/*.js; do
        if [ -f "$file" ]; then
            local basename=$(basename "$file")
            local version=$(echo "$basename" | cut -d'_' -f1)
            all_migrations+=("$version:$file")
        fi
    done
    
    # 排序
    IFS=$'\n' all_migrations=($(sort <<<"${all_migrations[*]}"))
    unset IFS
    
    # 找出待应用的迁移
    for migration in "${all_migrations[@]}"; do
        local version=$(echo "$migration" | cut -d':' -f1)
        local file=$(echo "$migration" | cut -d':' -f2)
        
        local is_applied=false
        for applied in "${applied_migrations[@]}"; do
            if [ "$version" = "$applied" ]; then
                is_applied=true
                break
            fi
        done
        
        if [ "$is_applied" = false ]; then
            echo "$file"
        fi
    done
}

# 应用单个迁移
apply_migration() {
    local file="$1"
    local uri="${MONGODB_URI:-$DEFAULT_MONGODB_URI}"
    local database="${DATABASE:-$DEFAULT_DATABASE}"
    
    log_info "应用迁移: $(basename "$file")"
    
    if command -v mongosh >/dev/null 2>&1; then
        if mongosh "$uri/$database" "$file"; then
            log_success "迁移应用成功: $(basename "$file")"
            return 0
        else
            log_error "迁移应用失败: $(basename "$file")"
            return 1
        fi
    else
        if mongo "$uri/$database" "$file"; then
            log_success "迁移应用成功: $(basename "$file")"
            return 0
        else
            log_error "迁移应用失败: $(basename "$file")"
            return 1
        fi
    fi
}

# 运行所有待应用的迁移
run_migrations() {
    log_info "开始数据库迁移..."
    
    local pending_migrations=()
    while IFS= read -r file; do
        if [ -n "$file" ]; then
            pending_migrations+=("$file")
        fi
    done < <(get_pending_migrations)
    
    if [ ${#pending_migrations[@]} -eq 0 ]; then
        log_info "没有待应用的迁移"
        return 0
    fi
    
    log_info "找到 ${#pending_migrations[@]} 个待应用的迁移"
    
    for file in "${pending_migrations[@]}"; do
        if ! apply_migration "$file"; then
            log_error "迁移失败，停止执行"
            return 1
        fi
    done
    
    log_success "所有迁移应用完成"
}

# 回滚迁移
rollback_migration() {
    local steps="${1:-1}"
    local uri="${MONGODB_URI:-$DEFAULT_MONGODB_URI}"
    local database="${DATABASE:-$DEFAULT_DATABASE}"
    
    log_info "回滚最近 $steps 个迁移..."
    
    # 获取最近应用的迁移
    local recent_migrations
    if command -v mongosh >/dev/null 2>&1; then
        recent_migrations=$(mongosh "$uri/$database" --quiet --eval "db.migrations.find({}).sort({applied_at: -1}).limit($steps).forEach(function(doc) { print(doc.filename); })")
    else
        recent_migrations=$(mongo "$uri/$database" --quiet --eval "db.migrations.find({}).sort({applied_at: -1}).limit($steps).forEach(function(doc) { print(doc.filename); })")
    fi
    
    if [ -z "$recent_migrations" ]; then
        log_warning "没有可回滚的迁移"
        return 0
    fi
    
    # 回滚每个迁移
    echo "$recent_migrations" | while IFS= read -r filename; do
        if [ -n "$filename" ]; then
            local filepath="$MIGRATIONS_DIR/$filename"
            if [ -f "$filepath" ]; then
                log_info "回滚迁移: $filename"
                
                # 创建临时回滚脚本
                local rollback_script="/tmp/rollback_$(basename "$filename")"
                
                # 提取回滚函数并执行
                cat > "$rollback_script" << EOF
$(cat "$filepath" | sed -n '/function down()/,/^}/p')

// 执行回滚
down();

// 删除迁移记录
db.migrations.deleteOne({filename: "$filename"});
EOF
                
                if command -v mongosh >/dev/null 2>&1; then
                    mongosh "$uri/$database" "$rollback_script"
                else
                    mongo "$uri/$database" "$rollback_script"
                fi
                
                rm "$rollback_script"
                log_success "回滚完成: $filename"
            else
                log_warning "迁移文件不存在: $filepath"
            fi
        fi
    done
    
    log_success "迁移回滚完成"
}

# 显示迁移状态
show_migration_status() {
    log_info "迁移状态:"
    
    # 显示已应用的迁移
    local applied_count=0
    echo "已应用的迁移:"
    while IFS= read -r version; do
        if [ -n "$version" ]; then
            echo "  ✓ $version"
            ((applied_count++))
        fi
    done < <(get_applied_migrations)
    
    # 显示待应用的迁移
    local pending_count=0
    echo "待应用的迁移:"
    while IFS= read -r file; do
        if [ -n "$file" ]; then
            echo "  ○ $(basename "$file")"
            ((pending_count++))
        fi
    done < <(get_pending_migrations)
    
    echo ""
    echo "统计: 已应用 $applied_count 个，待应用 $pending_count 个"
}

# 运行种子数据
run_seeds() {
    local uri="${MONGODB_URI:-$DEFAULT_MONGODB_URI}"
    local database="${DATABASE:-$DEFAULT_DATABASE}"
    
    log_info "运行种子数据..."
    
    if [ ! -d "$SEEDS_DIR" ]; then
        log_warning "种子数据目录不存在: $SEEDS_DIR"
        return 0
    fi
    
    local seed_files=("$SEEDS_DIR"/*.js)
    if [ ! -f "${seed_files[0]}" ]; then
        log_warning "没有找到种子数据文件"
        return 0
    fi
    
    for file in "${seed_files[@]}"; do
        if [ -f "$file" ]; then
            log_info "运行种子数据: $(basename "$file")"
            
            if command -v mongosh >/dev/null 2>&1; then
                mongosh "$uri/$database" "$file"
            else
                mongo "$uri/$database" "$file"
            fi
            
            log_success "种子数据运行完成: $(basename "$file")"
        fi
    done
    
    log_success "所有种子数据运行完成"
}

# 创建数据库索引
create_indexes() {
    local uri="${MONGODB_URI:-$DEFAULT_MONGODB_URI}"
    local database="${DATABASE:-$DEFAULT_DATABASE}"
    
    log_info "创建数据库索引..."
    
    # 创建索引脚本
    local index_script="/tmp/create_indexes.js"
    
    cat > "$index_script" << 'EOF'
// Greatest Works - 数据库索引创建脚本

print("开始创建索引...");

// 玩家集合索引
db.players.createIndex({ "user_id": 1 }, { unique: true });
db.players.createIndex({ "name": 1 }, { unique: true });
db.players.createIndex({ "level": -1 });
db.players.createIndex({ "created_at": 1 });
db.players.createIndex({ "last_login": -1 });
print("玩家集合索引创建完成");

// 公会集合索引
db.guilds.createIndex({ "name": 1 }, { unique: true });
db.guilds.createIndex({ "leader_id": 1 });
db.guilds.createIndex({ "level": -1 });
db.guilds.createIndex({ "created_at": 1 });
print("公会集合索引创建完成");

// 战斗记录索引
db.battles.createIndex({ "player_id": 1, "created_at": -1 });
db.battles.createIndex({ "battle_type": 1, "created_at": -1 });
db.battles.createIndex({ "created_at": -1 });
print("战斗记录索引创建完成");

// 排行榜索引
db.rankings.createIndex({ "type": 1, "score": -1 });
db.rankings.createIndex({ "player_id": 1, "type": 1 });
db.rankings.createIndex({ "updated_at": -1 });
print("排行榜索引创建完成");

// 聊天记录索引
db.chats.createIndex({ "channel_id": 1, "created_at": -1 });
db.chats.createIndex({ "player_id": 1, "created_at": -1 });
db.chats.createIndex({ "created_at": -1 }, { expireAfterSeconds: 2592000 }); // 30天过期
print("聊天记录索引创建完成");

// 日志集合索引
db.player_logs.createIndex({ "player_id": 1, "created_at": -1 });
db.player_logs.createIndex({ "action": 1, "created_at": -1 });
db.player_logs.createIndex({ "created_at": -1 }, { expireAfterSeconds: 7776000 }); // 90天过期
print("日志集合索引创建完成");

print("所有索引创建完成");
EOF
    
    if command -v mongosh >/dev/null 2>&1; then
        mongosh "$uri/$database" "$index_script"
    else
        mongo "$uri/$database" "$index_script"
    fi
    
    rm "$index_script"
    
    log_success "数据库索引创建完成"
}

# 备份数据库
backup_database() {
    local backup_dir="./backups"
    local timestamp=$(date +"%Y%m%d_%H%M%S")
    local backup_path="$backup_dir/backup_$timestamp"
    local uri="${MONGODB_URI:-$DEFAULT_MONGODB_URI}"
    local database="${DATABASE:-$DEFAULT_DATABASE}"
    
    log_info "备份数据库到: $backup_path"
    
    mkdir -p "$backup_dir"
    
    if command -v mongodump >/dev/null 2>&1; then
        mongodump --uri="$uri" --db="$database" --out="$backup_path"
        
        # 压缩备份
        tar -czf "$backup_path.tar.gz" -C "$backup_dir" "backup_$timestamp"
        rm -rf "$backup_path"
        
        log_success "数据库备份完成: $backup_path.tar.gz"
    else
        log_error "mongodump未安装，无法备份数据库"
        return 1
    fi
}

# 恢复数据库
restore_database() {
    local backup_file="$1"
    local uri="${MONGODB_URI:-$DEFAULT_MONGODB_URI}"
    local database="${DATABASE:-$DEFAULT_DATABASE}"
    
    if [ -z "$backup_file" ]; then
        log_error "请指定备份文件"
        return 1
    fi
    
    if [ ! -f "$backup_file" ]; then
        log_error "备份文件不存在: $backup_file"
        return 1
    fi
    
    log_info "从备份恢复数据库: $backup_file"
    
    # 解压备份文件
    local temp_dir="/tmp/restore_$(date +%s)"
    mkdir -p "$temp_dir"
    
    if [[ "$backup_file" == *.tar.gz ]]; then
        tar -xzf "$backup_file" -C "$temp_dir"
    else
        cp -r "$backup_file" "$temp_dir/"
    fi
    
    if command -v mongorestore >/dev/null 2>&1; then
        mongorestore --uri="$uri" --db="$database" --drop "$temp_dir"/*
        
        rm -rf "$temp_dir"
        
        log_success "数据库恢复完成"
    else
        log_error "mongorestore未安装，无法恢复数据库"
        rm -rf "$temp_dir"
        return 1
    fi
}

# 显示帮助信息
show_help() {
    echo "Greatest Works 数据库迁移脚本"
    echo ""
    echo "用法: $0 [命令] [选项]"
    echo ""
    echo "命令:"
    echo "  migrate                 运行所有待应用的迁移"
    echo "  rollback [N]            回滚最近N个迁移 (默认1个)"
    echo "  status                  显示迁移状态"
    echo "  generate NAME           生成新的迁移文件"
    echo "  seed                    运行种子数据"
    echo "  index                   创建数据库索引"
    echo "  backup                  备份数据库"
    echo "  restore FILE            从备份恢复数据库"
    echo "  setup                   初始化迁移环境"
    echo ""
    echo "选项:"
    echo "  -h, --help              显示帮助信息"
    echo "  --mongodb-uri URI       MongoDB连接URI"
    echo "  --database NAME         数据库名称"
    echo "  --redis-addr ADDR       Redis地址"
    echo ""
    echo "环境变量:"
    echo "  MONGODB_URI             MongoDB连接URI"
    echo "  DATABASE                数据库名称"
    echo "  REDIS_ADDR              Redis地址"
    echo ""
    echo "示例:"
    echo "  $0 migrate              # 运行迁移"
    echo "  $0 generate add_user_index  # 生成迁移文件"
    echo "  $0 rollback 2           # 回滚2个迁移"
    echo "  $0 backup               # 备份数据库"
}

# 主函数
main() {
    local command=""
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            --mongodb-uri)
                export MONGODB_URI="$2"
                shift 2
                ;;
            --database)
                export DATABASE="$2"
                shift 2
                ;;
            --redis-addr)
                export REDIS_ADDR="$2"
                shift 2
                ;;
            migrate|rollback|status|generate|seed|index|backup|restore|setup)
                command="$1"
                shift
                break
                ;;
            *)
                log_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    if [ -z "$command" ]; then
        log_error "请指定命令"
        show_help
        exit 1
    fi
    
    log_info "Greatest Works 数据库迁移脚本启动"
    log_info "命令: $command"
    
    # 执行命令
    case $command in
        "setup")
            setup_migration_dirs
            ;;
        "migrate")
            if check_mongodb; then
                setup_migration_dirs
                run_migrations
            fi
            ;;
        "rollback")
            if check_mongodb; then
                local steps="${1:-1}"
                rollback_migration "$steps"
            fi
            ;;
        "status")
            if check_mongodb; then
                show_migration_status
            fi
            ;;
        "generate")
            local name="$1"
            generate_migration "$name"
            ;;
        "seed")
            if check_mongodb; then
                run_seeds
            fi
            ;;
        "index")
            if check_mongodb; then
                create_indexes
            fi
            ;;
        "backup")
            if check_mongodb; then
                backup_database
            fi
            ;;
        "restore")
            if check_mongodb; then
                local backup_file="$1"
                restore_database "$backup_file"
            fi
            ;;
        *)
            log_error "未知命令: $command"
            show_help
            exit 1
            ;;
    esac
    
    log_success "操作完成"
}

# 执行主函数
main "$@"