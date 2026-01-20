# Push Notify Build Script (Windows PowerShell)

$ErrorActionPreference = "Stop"

$RootPath = Get-Location
$FrontendPath = Join-Path $RootPath "frontend"
$BackendPath = Join-Path $RootPath "backend"
$StaticDistPath = Join-Path $BackendPath "static\dist"

Write-Host ">>> Starting Build Process..." -ForegroundColor Cyan

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

# 3. Build Backend
Write-Host ">>> Building Backend..." -ForegroundColor Yellow
Set-Location $BackendPath
go build -o "$RootPath\push-notify.exe" main.go

Write-Host ">>> Build Finished! Binary is at: $RootPath\push-notify.exe" -ForegroundColor Green
Set-Location $RootPath
