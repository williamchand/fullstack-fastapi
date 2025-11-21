#!/bin/bash

# Build release script for multiple platforms

set -e

VERSION=${1:-$(git describe --tags --always)}
PLATFORMS="linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64"

echo "Building version: $VERSION"

# Create releases directory
mkdir -p releases

for platform in $PLATFORMS; do
    os=$(echo $platform | cut -d'/' -f1)
    arch=$(echo $platform | cut -d'/' -f2)
    
    output_name="releases/${BINARY_NAME}-${VERSION}-${os}-${arch}"
    
    if [ "$os" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    echo "Building for $os/$arch..."
    GOOS=$os GOARCH=$arch $(GOBUILD) -o $output_name ./cmd/server
    
    # Create archive
    if [ "$os" = "windows" ]; then
        zip -j "${output_name}.zip" $output_name
    else
        tar -czf "${output_name}.tar.gz" -C releases $(basename $output_name)
    fi
    
    rm $output_name
done

echo "Build complete! Release files in releases/"