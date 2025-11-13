@echo off
REM Proto文件生成脚本 (Windows版本)
REM 支持Go和C#代码生成

setlocal enabledelayedexpansion

echo 开始生成Proto文件...

REM 优先使用仓库自带的protoc
set PROTOC_PATH=%CD%\protoc\bin\protoc.exe
if exist "%PROTOC_PATH%" (
    set "PROTOC=%PROTOC_PATH%"
) else (
    where protoc >nul 2>nul
    if %errorlevel% equ 0 (
        for /f "delims=" %%i in ('where protoc') do set "PROTOC=%%i"
    ) else (
        echo 错误: 未找到protoc，请安装或将其添加到PATH，或将protoc放置在repo的protoc\bin目录下。
        exit /b 1
    )
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
if not exist "internal\proto\chat" mkdir "internal\proto\chat"
if not exist "internal\proto\team" mkdir "internal\proto\team"
if not exist "internal\proto\mail" mkdir "internal\proto\mail"
if not exist "internal\proto\room" mkdir "internal\proto\room"
if not exist "internal\proto\scene" mkdir "internal\proto\scene"
if not exist "internal\proto\gateway" mkdir "internal\proto\gateway"
if not exist "csharp\GreatestWorks\Player" mkdir "csharp\GreatestWorks\Player"
if not exist "csharp\GreatestWorks\Battle" mkdir "csharp\GreatestWorks\Battle"
if not exist "csharp\GreatestWorks\Pet" mkdir "csharp\GreatestWorks\Pet"
if not exist "csharp\GreatestWorks\Common" mkdir "csharp\GreatestWorks\Common"
if not exist "csharp\GreatestWorks\Chat" mkdir "csharp\GreatestWorks\Chat"
if not exist "csharp\GreatestWorks\Team" mkdir "csharp\GreatestWorks\Team"
if not exist "csharp\GreatestWorks\Mail" mkdir "csharp\GreatestWorks\Mail"
if not exist "csharp\GreatestWorks\Room" mkdir "csharp\GreatestWorks\Room"
if not exist "csharp\GreatestWorks\Scene" mkdir "csharp\GreatestWorks\Scene"
if not exist "csharp\GreatestWorks\Gateway" mkdir "csharp\GreatestWorks\Gateway"

echo 使用的protoc: %PROTOC%
echo 生成Go代码...

REM 生成Go代码
"%PROTOC%" --go_out=. --go_opt=paths=source_relative proto\common.proto
"%PROTOC%" --go_out=. --go_opt=paths=source_relative proto\player.proto
"%PROTOC%" --go_out=. --go_opt=paths=source_relative proto\battle.proto
"%PROTOC%" --go_out=. --go_opt=paths=source_relative proto\pet.proto
"%PROTOC%" --go_out=. --go_opt=paths=source_relative proto\chat.proto
"%PROTOC%" --go_out=. --go_opt=paths=source_relative proto\team.proto
"%PROTOC%" --go_out=. --go_opt=paths=source_relative proto\mail.proto
"%PROTOC%" --go_out=. --go_opt=paths=source_relative proto\room.proto
"%PROTOC%" --go_out=. --go_opt=paths=source_relative proto\scene.proto
"%PROTOC%" --go_out=. --go_opt=paths=source_relative proto\gateway.proto

echo 移动生成的文件到正确位置...

REM 移动生成的文件到正确位置
move proto\common.pb.go internal\proto\common\ >nul 2>nul
move proto\player.pb.go internal\proto\player\ >nul 2>nul
move proto\battle.pb.go internal\proto\battle\ >nul 2>nul
move proto\pet.pb.go internal\proto\pet\ >nul 2>nul
move proto\chat.pb.go internal\proto\chat\ >nul 2>nul
move proto\team.pb.go internal\proto\team\ >nul 2>nul
move proto\mail.pb.go internal\proto\mail\ >nul 2>nul
move proto\room.pb.go internal\proto\room\ >nul 2>nul
move proto\scene.pb.go internal\proto\scene\ >nul 2>nul
move proto\gateway.pb.go internal\proto\gateway\ >nul 2>nul



echo 生成C#代码...

REM 生成C#代码
"%PROTOC%" --csharp_out=csharp --csharp_opt=file_extension=.g.cs proto\common.proto
"%PROTOC%" --csharp_out=csharp --csharp_opt=file_extension=.g.cs proto\player.proto
"%PROTOC%" --csharp_out=csharp --csharp_opt=file_extension=.g.cs proto\battle.proto
"%PROTOC%" --csharp_out=csharp --csharp_opt=file_extension=.g.cs proto\pet.proto
"%PROTOC%" --csharp_out=csharp --csharp_opt=file_extension=.g.cs proto\chat.proto
"%PROTOC%" --csharp_out=csharp --csharp_opt=file_extension=.g.cs proto\team.proto
"%PROTOC%" --csharp_out=csharp --csharp_opt=file_extension=.g.cs proto\mail.proto
"%PROTOC%" --csharp_out=csharp --csharp_opt=file_extension=.g.cs proto\room.proto
"%PROTOC%" --csharp_out=csharp --csharp_opt=file_extension=.g.cs proto\scene.proto
"%PROTOC%" --csharp_out=csharp --csharp_opt=file_extension=.g.cs proto\gateway.proto

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
