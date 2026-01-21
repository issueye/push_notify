# Push Notify Cross-Compile Script (Windows to Linux)
$ErrorActionPreference = "Stop"

$RootPath = Get-Location
$FrontendPath = Join-Path $RootPath "frontend"
$BackendPath = Join-Path $RootPath "backend"
$StaticDistPath = Join-Path $BackendPath "static\dist"

Write-Host ">>> Starting Cross-Compile Process (Windows -> Linux)..." -ForegroundColor Cyan

# 1. Build Frontend
Write-Host ">>> Building Frontend..." -ForegroundColor Yellow
Set-Location $FrontendPath
npm install
npm run build

# 2. Prepare Backend Static Directory
Write-Host ">>> Preparing Backend Static Assets..." -ForegroundColor Yellow
if (Test-Path $StaticDistPath) {
    Remove-Item -Recurse -Force $StaticDistPath
}
New-Item -ItemType Directory -Force -Path $StaticDistPath | Out-Null

# Copy dist content to backend/static/dist
Copy-Item -Recurse -Force "$FrontendPath\dist\*" $StaticDistPath

# 3. Build Backend for Linux
Write-Host ">>> Building Backend for Linux (AMD64)..." -ForegroundColor Yellow
Set-Location $BackendPath
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go build -o "$RootPath\push-notify-linux" main.go
# Reset environment variables
$env:GOOS = ""
$env:GOARCH = ""

Write-Host ">>> Build Finished! Linux binary is at: $RootPath\push-notify-linux" -ForegroundColor Green
Set-Location $RootPath
