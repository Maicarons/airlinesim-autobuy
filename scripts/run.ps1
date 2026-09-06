#!/usr/bin/env pwsh
# Run script for AirlineSim Autobuy (Windows PowerShell)
# Usage: .\scripts\run.ps1 [config_path]

param(
    [string]$ConfigPath = "configs/config.yaml"
)

if (-not (Test-Path $ConfigPath)) {
    Write-Host "Error: Config file not found: $ConfigPath"
    Write-Host "Usage: .\scripts\run.ps1 [[-ConfigPath] <path>]"
    Write-Host ""
    Write-Host "If you haven't set up the config yet, run:"
    Write-Host "  Copy-Item configs/config.example.yaml configs/config.yaml"
    Write-Host "  # Then edit configs/config.yaml with your credentials"
    exit 1
}

Write-Host "Starting AirlineSim Autobuy with config: $ConfigPath"
Write-Host "Web UI will be available at http://localhost:9090"
Write-Host ""

if (-not (Test-Path autobuy.exe)) {
    Write-Host "Binary not found. Building first..."
    & make build
    Write-Host ""
}

& .\autobuy.exe