[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
$OutputEncoding = [System.Text.Encoding]::UTF8

$ErrorActionPreference = "Stop"

Write-Host ""
Write-Host "========================================"
Write-Host "POC Vision API - Build & Deploy"
Write-Host "========================================"
Write-Host ""

$rootDir = Split-Path -Parent $PSScriptRoot
Push-Location $rootDir

# Build
Write-Host "Step 1/2: Building lambdas..."
& "$PSScriptRoot\build.ps1"
if ($LASTEXITCODE -ne 0) {
    Write-Error "Build failed"
    Pop-Location
    exit 1
}

Write-Host ""
Write-Host "Step 2/2: Deploying to AWS..."

# Deploy
& "$PSScriptRoot\deploy.ps1"
if ($LASTEXITCODE -ne 0) {
    Write-Error "Deploy failed"
    Pop-Location
    exit 1
}

Pop-Location

Write-Host ""
Write-Host "========================================"
Write-Host "Build and deploy completed successfully!"
Write-Host "========================================"
