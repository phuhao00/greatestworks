@echo off
REM 启动分布式游戏服务脚本
REM 基于DDD架构的分布式多节点服务

echo 启动分布式游戏服务...

REM 检查Go环境
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo 错误: 未找到Go环境，请先安装Go
    pause
    exit /b 1
)

REM 创建日志目录
if not exist "logs" mkdir logs

REM 启动认证服务
echo 启动认证服务...
start "Auth Service" cmd /k "cd /d %~dp0.. && go run cmd/auth-service/main.go"

REM 等待认证服务启动
timeout /t 3 /nobreak >nul

REM 启动游戏服务
echo 启动游戏服务...
start "Game Service" cmd /k "cd /d %~dp0.. && go run cmd/game-service/main.go"

REM 等待游戏服务启动
timeout /t 3 /nobreak >nul

REM 启动网关服务
echo 启动网关服务...
start "Gateway Service" cmd /k "cd /d %~dp0.. && go run cmd/gateway-service/main.go"

echo.
echo 所有服务已启动！
echo.
echo 服务地址：
echo - 认证服务: http://localhost:8080
echo - 游戏服务: rpc://localhost:8081
echo - 网关服务: tcp://localhost:9090
echo.
echo 按任意键关闭所有服务...

pause >nul

REM 关闭所有服务
echo 正在关闭所有服务...
taskkill /f /im go.exe >nul 2>&1
taskkill /f /im cmd.exe >nul 2>&1

echo 所有服务已关闭
pause
