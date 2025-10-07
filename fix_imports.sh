#!/bin/bash

# 批量修复导入路径
find . -name "*.go" -type f -exec sed -i 's|greatestworks/internal/infrastructure/logger|greatestworks/internal/infrastructure/logging|g' {} \;

echo "导入路径修复完成"
