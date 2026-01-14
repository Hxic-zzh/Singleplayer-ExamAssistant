@echo off
chcp 65001 > nul
echo ========================================
echo  Wails 开发调试启动器
echo ========================================
echo.

REM 检查是否安装了 wails
where wails >nul 2>nul
if %errorlevel% neq 0 (
    echo [错误] 未找到 wails 命令
    echo 请先运行: go install github.com/wailsapp/wails/v2/cmd/wails@latest
    pause
    exit /b 1
)

REM 检查是否在项目目录中
if not exist "wails.json" (
    echo [错误] 不在 Wails 项目目录中
    echo 请确保当前目录包含 wails.json 文件
    pause
    exit /b 1
)

echo [信息] 检查依赖...
go mod tidy

echo [信息] 清理之前的构建...
if exist "build" (
    rmdir /s /q build 2>nul
    echo 已清理 build 目录
)

echo [信息] 启动 Wails 开发服务器...
echo ========================================
echo 开发服务器启动中...
echo 前端地址: http://localhost:34115
echo 热重载: 已启用
echo 日志输出: 已启用
echo 按 Ctrl+C 停止服务器
echo ========================================
echo.

wails dev

if %errorlevel% neq 0 (
    echo.
    echo [错误] 启动失败，错误代码: %errorlevel%
    echo 可能的原因:
    echo 1. 端口被占用
    echo 2. 依赖缺失
    echo 3. 前端文件错误
    pause
)