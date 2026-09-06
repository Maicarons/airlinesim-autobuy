#!/usr/bin/env pwsh
# Install script for AirlineSim Autobuy (Windows PowerShell)
# Usage: .\scripts\install.ps1

Write-Host "============================================"
Write-Host " AirlineSim Autobuy - Install"
Write-Host "============================================"
Write-Host ""

# Check prerequisites
Write-Host "Checking prerequisites..."

$hasGo = $null -ne (Get-Command go -ErrorAction SilentlyContinue)
if ($hasGo) {
    Write-Host "  ✓ Go found"
} else {
    Write-Host "  ✗ Go not found - please install Go 1.22+ from https://go.dev"
    exit 1
}

$hasNode = $null -ne (Get-Command node -ErrorAction SilentlyContinue)
if ($hasNode) {
    Write-Host "  ✓ Node.js found"
} else {
    Write-Host "  ⚠ Node.js not found - skipping frontend build"
}

Write-Host ""

# Build frontend
if ($hasNode) {
    Write-Host "Installing frontend dependencies..."
    Push-Location internal/webui/frontend
    npm install
    Write-Host "Building frontend..."
    npm run build
    Pop-Location
    Write-Host "  ✓ Frontend built"
}

# Build backend
Write-Host "Building Go backend..."
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -ldflags="-s -w" -o autobuy.exe ./cmd/autobuy/
Write-Host "  ✓ Backend built: ./autobuy.exe"

Write-Host ""

# Setup config
if (-not (Test-Path configs/config.yaml)) {
    Write-Host "Creating default config..."
    Copy-Item configs/config.example.yaml configs/config.yaml
    Write-Host "  ✓ Created configs/config.yaml"
    Write-Host "  ⚠ Edit configs/config.yaml with your credentials before running"
} else {
    Write-Host "  ✓ configs/config.yaml already exists"
}

Write-Host ""
Write-Host "============================================"
Write-Host " Installation complete!"
Write-Host ""
Write-Host " Next steps:"
Write-Host "   1. Edit configs/config.yaml with your credentials"
Write-Host "   2. Run: .\autobuy.exe"
Write-Host "   3. Open http://localhost:9090"
Write-Host "============================================"