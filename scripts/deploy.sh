#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

if [ -f "./config.env" ]; then
    echo "Loading environment variables from config.env"
    export $(grep -v '^#' config.env | xargs)
fi

ENVIRONMENT=${ENVIRONMENT:-dev}
CATALOG_API_TOKEN=${CATALOG_API_TOKEN:-poc-vision-api-token-2026}
BEDROCK_MODEL_ID=${BEDROCK_MODEL_ID:-anthropic.claude-sonnet-4-20250514-v1:0}

STACK_NAME="poc-vision-api-${ENVIRONMENT}"

required_zips=("main-lambda.zip" "mcp-server.zip" "catalog-api.zip")
for zip in "${required_zips[@]}"; do
    if [ ! -f "build/$zip" ]; then
        echo "Build artifacts not found. Running build script..."
        "$SCRIPT_DIR/build.sh"
        break
    fi
done

echo "Deploying SAM stack: $STACK_NAME"

sam deploy \
    --template-file deployments/template.yaml \
    --stack-name "$STACK_NAME" \
    --capabilities CAPABILITY_IAM \
    --parameter-overrides \
        Environment="$ENVIRONMENT" \
        CatalogApiToken="$CATALOG_API_TOKEN" \
        BedrockModelId="$BEDROCK_MODEL_ID" \
    --resolve-s3 \
    --no-confirm-changeset

if [ $? -ne 0 ]; then
    echo "SAM deployment failed"
    exit 1
fi

echo "Getting stack outputs..."

ANALYZE_ENDPOINT=$(aws cloudformation describe-stacks \
    --stack-name "$STACK_NAME" \
    --query "Stacks[0].Outputs[?OutputKey=='AnalyzeEndpoint'].OutputValue" \
    --output text)

CATALOG_ENDPOINT=$(aws cloudformation describe-stacks \
    --stack-name "$STACK_NAME" \
    --query "Stacks[0].Outputs[?OutputKey=='CatalogEndpoint'].OutputValue" \
    --output text)

API_KEY=$(aws cloudformation describe-stacks \
    --stack-name "$STACK_NAME" \
    --query "Stacks[0].Outputs[?OutputKey=='ApiKeyValue'].OutputValue" \
    --output text)

MODEL_USED=$(aws cloudformation describe-stacks \
    --stack-name "$STACK_NAME" \
    --query "Stacks[0].Outputs[?OutputKey=='BedrockModelUsed'].OutputValue" \
    --output text)

echo ""
echo "Deployment completed successfully"
echo "=================================="
echo "Analyze Endpoint: $ANALYZE_ENDPOINT"
echo "Catalog Endpoint: $CATALOG_ENDPOINT"
echo "API Key: $API_KEY"
echo "Bedrock Model: $MODEL_USED"
echo ""
echo "Test command:"
echo "curl -X POST $ANALYZE_ENDPOINT -H 'x-api-key: $API_KEY' -H 'Content-Type: application/json' -d '{\"image_base64\":\"<base64>\",\"media_type\":\"image/jpeg\"}'"
