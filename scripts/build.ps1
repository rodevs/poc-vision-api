[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
$OutputEncoding = [System.Text.Encoding]::UTF8

$ErrorActionPreference = "Stop"

Write-Host ""
Write-Host "========================================"
Write-Host "POC Vision API - Build"
Write-Host "========================================"
Write-Host ""

$rootDir = Split-Path -Parent $PSScriptRoot
Push-Location $rootDir

$buildDir = Join-Path -Path $rootDir -ChildPath "build"

if (-not (Test-Path $buildDir)) {
    New-Item -ItemType Directory -Path $buildDir | Out-Null
}

# Tidy dependencies
Write-Host "Running go mod tidy..."
go mod tidy
if ($LASTEXITCODE -ne 0) {
    Write-Error "go mod tidy failed"
    Pop-Location
    Exit 1
}

# Lambda configurations
$lambdas = @(
    @{ Name = "main-lambda"; Path = "./cmd/main-lambda" },
    @{ Name = "mcp-server"; Path = "./cmd/mcp-server" },
    @{ Name = "catalog-api"; Path = "./cmd/catalog-api" }
)

$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

foreach ($lambda in $lambdas) {
    Write-Host "Building $($lambda.Name)..."

    $outputPath = Join-Path -Path $buildDir -ChildPath "bootstrap"
    
    go build -tags lambda.norpc -o $outputPath $lambda.Path
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Build failed for $($lambda.Name)"
        Pop-Location
        Exit 1
    }

    $zipFile = Join-Path -Path $buildDir -ChildPath "$($lambda.Name).zip"
    
    if (Test-Path $zipFile) {
        Remove-Item $zipFile -Force
    }

    Push-Location $buildDir
    Compress-Archive -Path "bootstrap" -DestinationPath $zipFile -CompressionLevel Optimal
    Pop-Location

    Remove-Item $outputPath -Force

    Write-Host "  -> $($lambda.Name).zip created"
}

Pop-Location

Write-Host ""
Write-Host "========================================"
Write-Host "Build completed successfully!"
Write-Host "========================================"
Write-Host "Artifacts in: $buildDir"
Write-Host ""
