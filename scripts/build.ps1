#!/usr/bin/env pwsh
# Build script for AirlineSim Autobuy (Windows PowerShell)
# Usage: .\scripts\build.ps1 [linux|windows|darwin|all]

$AppName = "autobuy"
$BuildDir = "build"
$Version = if (git describe --tags --always 2>$null) { git describe --tags --always } else { "dev" }
$Commit = if (git rev-parse --short HEAD 2>$null) { git rev-parse --short HEAD } else { "unknown" }
$Date = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")
$LdFlags = "-s -w -X main.version=$Version -X main.commit=$Commit -X main.date=$Date"

New-Item -ItemType Directory -Force -Path $BuildDir | Out-Null

function Build-Platform($goos, $goarch, $ext) {
    $output = "$BuildDir/$AppName-$goos-$goarch$ext"
    Write-Host "Building for $goos/$goarch..."
    $env:GOOS = $goos
    $env:GOARCH = $goarch
    $env:CGO_ENABLED = "0"
    & go build -ldflags=$LdFlags -o $output ./cmd/autobuy/
    Write-Host "  → $output"
}

switch ($args[0]) {
    "linux" {
        Build-Platform "linux" "amd64" ""
        Build-Platform "linux" "arm64" ""
    }
    "windows" {
        Build-Platform "windows" "amd64" ".exe"
        Build-Platform "windows" "arm64" ".exe"
    }
    "darwin" {
        Build-Platform "darwin" "amd64" ""
        Build-Platform "darwin" "arm64" ""
    }
    default {
        Build-Platform "linux" "amd64" ""
        Build-Platform "linux" "arm64" ""
        Build-Platform "windows" "amd64" ".exe"
        Build-Platform "darwin" "amd64" ""
        Build-Platform "darwin" "arm64" ""
    }
}

Write-Host "`nBuild complete! Binaries in $BuildDir/"