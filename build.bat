@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

:: 设置控制台标题
title Hxic 应用程序编译工具

echo 正在编译 Hxic 应用程序...

:: 创建输出目录
if not exist "dist" mkdir dist

:: Windows 64位 - 带控制台版本
echo 编译 Windows 64位版本（带控制台）...
set GOOS=windows
set GOARCH=amd64
go build -o dist/Hxic-Windows64-Console.exe main.go
if !errorlevel! neq 0 (
    echo 错误：64位控制台版本编译失败！
    goto :error
)

:: Windows 64位 - 无控制台版本
echo 编译 Windows 64位版本（无控制台）...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-H windowsgui" -o dist/Hxic-Windows64-GUI.exe main.go
if !errorlevel! neq 0 (
    echo 错误：64位GUI版本编译失败！
    goto :error
)

:: 注意：Fyne 不支持 32位 Windows 编译
echo 跳过 32位版本编译（Fyne 不支持 32位 Windows）

:: 检查文件是否生成成功
set file_count=0
for /f %%i in ('dir /b dist\*.exe 2^>nul') do set /a file_count+=1

if !file_count! lss 2 (
    echo 警告：只成功生成了 !file_count! 个文件，预期 2 个文件
) else (
    echo 成功生成了所有 2 个 64位版本
)

echo.
echo ========================================
echo           编译完成！
echo ========================================
echo.
echo 生成的文件（在 dist 目录中）：
echo.
echo   ■ Hxic-Windows64-Console.exe
echo     64位带控制台版本（用于调试）
echo.
echo   ■ Hxic-Windows64-GUI.exe
echo     64位无控制台版本（用于正式发布）
echo.
echo ========================================
echo.
echo 文件大小统计：
dir dist\*.exe
echo.
goto :success

:error
echo.
echo ========================================
echo           编译过程中出现错误！
echo ========================================
echo 请检查：
echo   1. Go 环境是否配置正确
echo   2. 项目依赖是否完整
echo   3. 代码是否有语法错误
echo   4. 注意：Fyne 不支持 32位 Windows 编译
echo.
pause
exit /b 1

:success
echo 按任意键退出...
pause >nul
exit /b 0