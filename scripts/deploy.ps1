[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
$OutputEncoding = [System.Text.Encoding]::UTF8

$ErrorActionPreference = "Stop"

# Configuración por defecto
$DEFAULT_STACK_NAME = "poc-vision-api"
$DEFAULT_AI_PROVIDER = "bedrock"
$DEFAULT_BEDROCK_MODEL_ID = "amazon.nova-lite-v1:0"
$DEFAULT_AWS_REGION = "us-east-1"
$DEFAULT_CATALOG_API_TOKEN = "poc-vision-api-token-2026"

function Load-EnvFile {
    param([string]$filePath)
    if (Test-Path $filePath) {
        Write-Host "Loading environment variables from $filePath"
        Get-Content $filePath | ForEach-Object {
            if ($_ -match '^\s*([^#][^=]*)=(.*)$') {
                $name = $matches[1].Trim()
                $value = $matches[2].Trim()
                if (-not (Get-ChildItem Env: | Where-Object { $_.Name -eq $name })) {
                    Set-Item -Path ("Env:" + $name) -Value $value
                }
            }
        }
    }
}

# Cargar variables desde archivo si existe
$rootDir = Split-Path -Parent $PSScriptRoot
Load-EnvFile -filePath (Join-Path $rootDir "config.env")

# Aplicar valores por defecto
if (-not $env:AWS_REGION) { $env:AWS_REGION = $DEFAULT_AWS_REGION }
if (-not $env:AI_PROVIDER) { $env:AI_PROVIDER = $DEFAULT_AI_PROVIDER }
if (-not $env:BEDROCK_MODEL_ID) { $env:BEDROCK_MODEL_ID = $DEFAULT_BEDROCK_MODEL_ID }
if (-not $env:CATALOG_API_TOKEN) { $env:CATALOG_API_TOKEN = $DEFAULT_CATALOG_API_TOKEN }
if (-not $env:AI_API_KEY) { $env:AI_API_KEY = "" }

$stackName = $DEFAULT_STACK_NAME

$requiredZips = @("main-lambda.zip", "mcp-server.zip", "catalog-api.zip")
$buildDir = Join-Path -Path $rootDir -ChildPath "build"

foreach ($zip in $requiredZips) {
    $zipPath = Join-Path -Path $buildDir -ChildPath $zip
    if (-not (Test-Path $zipPath)) {
        Write-Host "Build artifacts not found. Running build script..."
        & "$PSScriptRoot\build.ps1"
        break
    }
}

$paramOverrides = @(
    "AIProvider=$($env:AI_PROVIDER)"
    "BedrockModelId=$($env:BEDROCK_MODEL_ID)"
    "CatalogApiToken=$($env:CATALOG_API_TOKEN)"
    "AIApiKey=$($env:AI_API_KEY)"
) -join " "

Write-Host ""
Write-Host "========================================"
Write-Host "POC Vision API - Deployment"
Write-Host "========================================"
Write-Host "Stack Name:   $stackName"
Write-Host "Region:       $($env:AWS_REGION)"
Write-Host "AI Provider:  $($env:AI_PROVIDER)"
Write-Host "Model:        $($env:BEDROCK_MODEL_ID)"
Write-Host "========================================"
Write-Host ""

Push-Location $rootDir

sam deploy `
    --template-file deployments/template.yaml `
    --stack-name $stackName `
    --capabilities CAPABILITY_IAM `
    --parameter-overrides $paramOverrides `
    --resolve-s3 `
    --no-confirm-changeset `
    --region $env:AWS_REGION

if ($LASTEXITCODE -ne 0) {
    Pop-Location
    Write-Error "SAM deployment failed"
    exit 1
}

Write-Host "Getting stack outputs..."

try {
    $outputs = aws cloudformation describe-stacks --stack-name $stackName --region $env:AWS_REGION --query "Stacks[0].Outputs" --output json | ConvertFrom-Json

    if (-not $outputs) {
        Write-Error "Failed to get stack outputs"
        exit 1
    }
} catch {
    Write-Error "Error getting stack outputs: $_"
    exit 1
}

$analyzeEndpoint = ($outputs | Where-Object { $_.OutputKey -eq "AnalyzeEndpoint" }).OutputValue
$catalogEndpoint = ($outputs | Where-Object { $_.OutputKey -eq "CatalogEndpoint" }).OutputValue
$apiKey = ($outputs | Where-Object { $_.OutputKey -eq "ApiKeyValue" }).OutputValue
$modelUsed = ($outputs | Where-Object { $_.OutputKey -eq "BedrockModelUsed" }).OutputValue

Write-Host ""
Write-Host "========================================"
Write-Host "Deployment completed successfully!"
Write-Host "========================================"
Write-Host "Analyze Endpoint: $analyzeEndpoint"
Write-Host "Catalog Endpoint: $catalogEndpoint"
Write-Host "API Key:          $apiKey"
Write-Host "Bedrock Model:    $modelUsed"
Write-Host "========================================"
Write-Host ""
Write-Host "Test catalog:"
Write-Host "  curl -H 'x-api-key: $apiKey' $catalogEndpoint"
Write-Host ""
Write-Host "Test analyze:"
Write-Host "  curl -X POST $analyzeEndpoint -H 'x-api-key: $apiKey' -H 'Content-Type: application/json' -d '{\"image_base64\":\"<base64>\",\"media_type\":\"image/jpeg\"}'"

Pop-Location
