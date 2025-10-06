@echo off
REM Proto文件生成脚本 (Windows版本)
REM 支持Go和C#代码生成

setlocal enabledelayedexpansion

echo 开始生成Proto文件...

REM 检查protoc是否安装
where protoc >nul 2>nul
if %errorlevel% neq 0 (
    echo 错误: protoc未安装，正在下载安装...
    echo 下载protoc...
    powershell -Command "Invoke-WebRequest -Uri 'https://github.com/protocolbuffers/protobuf/releases/download/v25.1/protoc-25.1-win64.zip' -OutFile 'protoc.zip'"
    echo 解压protoc...
    powershell -Command "Expand-Archive -Path 'protoc.zip' -DestinationPath 'protoc' -Force"
    echo 设置环境变量...
    set PATH=%PATH%;%CD%\protoc\bin
    echo protoc安装完成
)

REM 检查Go插件
where protoc-gen-go >nul 2>nul
if %errorlevel% neq 0 (
    echo 安装Go插件...
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
)

REM 创建输出目录
if not exist "internal\proto\player" mkdir "internal\proto\player"
if not exist "internal\proto\battle" mkdir "internal\proto\battle"
if not exist "internal\proto\pet" mkdir "internal\proto\pet"
if not exist "internal\proto\common" mkdir "internal\proto\common"
if not exist "csharp\GreatestWorks\Player" mkdir "csharp\GreatestWorks\Player"
if not exist "csharp\GreatestWorks\Battle" mkdir "csharp\GreatestWorks\Battle"
if not exist "csharp\GreatestWorks\Pet" mkdir "csharp\GreatestWorks\Pet"
if not exist "csharp\GreatestWorks\Common" mkdir "csharp\GreatestWorks\Common"

echo 生成Go代码...

REM 生成Go代码
protoc --go_out=. --go_opt=paths=source_relative proto\common.proto
protoc --go_out=. --go_opt=paths=source_relative proto\player.proto
protoc --go_out=. --go_opt=paths=source_relative proto\battle.proto
protoc --go_out=. --go_opt=paths=source_relative proto\pet.proto

echo 移动生成的文件到正确位置...

REM 移动生成的文件到正确位置
move proto\common.pb.go internal\proto\common\ >nul 2>nul
move proto\player.pb.go internal\proto\player\ >nul 2>nul
move proto\battle.pb.go internal\proto\battle\ >nul 2>nul
move proto\pet.pb.go internal\proto\pet\ >nul 2>nul



echo 生成C#代码...

REM 生成C#代码
protoc --csharp_out=csharp --csharp_opt=file_extension=.g.cs proto\common.proto
protoc --csharp_out=csharp --csharp_opt=file_extension=.g.cs proto\player.proto
protoc --csharp_out=csharp --csharp_opt=file_extension=.g.cs proto\battle.proto
protoc --csharp_out=csharp --csharp_opt=file_extension=.g.cs proto\pet.proto

echo Proto文件生成完成！
echo 生成的文件:
echo   Go代码: internal\proto\
echo   C#代码: csharp\

REM 显示生成的文件
echo Go生成的文件:
dir /s /b internal\proto\*.pb.go

echo C#生成的文件:
dir /s /b csharp\*.g.cs

echo.


pause
