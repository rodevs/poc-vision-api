#!/bin/bash

set -e

echo "Building POC Vision API Lambdas for Linux..."

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
BUILD_DIR="$PROJECT_DIR/build"

mkdir -p "$BUILD_DIR"

paths=("cmd/main-lambda" "cmd/mcp-server" "cmd/catalog-api")

for path in "${paths[@]}"; do
    echo "Building $path ..."

    cd "$PROJECT_DIR/$path"

    if [ ! -f "go.mod" ]; then
        echo "Error: go.mod not found in $path"
        exit 1
    fi

    go mod tidy

    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap main.go

    zipFileName="${path#cmd/}.zip"
    zipFile="$BUILD_DIR/$zipFileName"

    rm -f "$zipFile"
    zip "$zipFile" bootstrap
    rm bootstrap

    echo "$path built and packaged: $zipFile"
done

echo "Build completed successfully."
