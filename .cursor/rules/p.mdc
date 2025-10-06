---
description: Go项目开发规则 - 遵循Go最佳实践，不使用checkout和reset命令
globs: ["**/*.go", "**/go.mod", "**/go.sum", "**/*.yaml", "**/*.yml"]
alwaysApply: true
---

# Go项目开发规则

## 代码规范
- 遵循Go官方代码规范
- 使用gofmt格式化代码
- 变量命名使用驼峰命名法
- 常量使用大写字母和下划线
- 包名使用小写字母，简短有意义

## 错误处理
- 始终检查错误返回值
- 使用errors.New()或fmt.Errorf()创建错误
- 避免使用panic，除非是程序启动时的致命错误
- 使用defer进行资源清理

## 并发编程
- 使用goroutine处理并发任务
- 使用channel进行goroutine间通信
- 避免共享内存，使用通信来共享内存
- 使用sync包进行同步控制

## 性能优化
- 避免不必要的内存分配
- 使用对象池复用对象
- 合理使用缓存
- 避免在循环中进行重复计算


## 项目结构
- 遵循标准Go项目布局
- 按功能模块组织代码
- 使用internal包限制内部代码访问
- 合理使用接口和依赖注入

## 禁止使用的命令
- 不使用git checkout命令
- 不使用git reset命令
- 使用git switch切换分支
- 使用git restore恢复文件
