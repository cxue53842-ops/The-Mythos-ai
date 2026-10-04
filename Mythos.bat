@echo off
title THE MYTHOS // COMPILATION CORE v5.1
color 0A
echo =======================================================
echo [⚠] INITIALIZING PHYSICAL BUILD PIPELINE
echo =======================================================

:: 1. 檢測物理環境依賴
echo [*] Checking local environment dependencies...
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo [CRITICAL ERROR] Go compiler NOT found in system PATH.
    echo Please download from https://go.dev and retry.
    pause
    exit /b
)
where git >nul 2>nul
if %errorlevel% neq 0 (
    echo [CRITICAL ERROR] Git client NOT found in system PATH.
    pause
    exit /b
)
where gcc >nul 2>nul
if %errorlevel% neq 0 (
    echo [WARN] GCC compiler not detected. Fyne GUI graphics require MinGW-w64.
    echo Attempting to proceed, but compilation might drop.
)

:: 2. 建立臨時沙盒隔離區並從 GitHub 下載
echo [*] Creating isolated compile workspace...
if exist "mythos_build_temp" rmdir /s /q "mythos_build_temp"
mkdir "mythos_build_temp"
cd "mythos_build_temp"

echo [*] Downloading real source code matrix from GitHub...
:: [真實網絡拉取] 請將下方 URL 替換為你真實的 GitHub 專案倉庫路徑
git clone https://github.com/cxue53842-ops/The-Mythos-ai/tree/main .
if %errorlevel% neq 0 (
    echo [CRITICAL ERROR] Failed to fetch source from remote repository network node.
    cd ..
    pause
    exit /b
)

:: 3. 模組治理與依賴整合
echo [*] Fetching upstream Fyne GUI and cryptographic structures...
go mod tidy
if %errorlevel% neq 0 (
    echo [*] Mod file missing. Initializing new memory module matrix...
    go mod init themythos
    go get fyne.io/fyne/v2
    go mod tidy
)

:: 4. 執行黑客無痕編譯封裝
echo [*] Executing physical compilation sequence...
echo [*] Stripping OS terminal window attachments (-H=windowsgui)...
go build -ldflags="-H=windowsgui -s -w" -o "The_Mythos_Core.exe" main.go
if %errorlevel% neq 0 (
    echo [CRITICAL ERROR] Compilation fault. Syntax matrix or CGO linkage rejected.
    cd ..
    pause
    exit /b
)

:: 5. 數據遷移與清理
echo [*] Extracting final binary output...
move "The_Mythos_Core.exe" "../The_Mythos_Core.exe" >nul
cd ..
rmdir /s /q "mythos_build_temp"

echo =======================================================
echo [SUCCESS] THE MYTHOS BINARY CORE DEPLOYED SUCCESSFULLY
echo TARGET OUTPUT: ./The_Mythos_Core.exe
echo =======================================================
pause
