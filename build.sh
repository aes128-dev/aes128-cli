#!/bin/bash
set -e

VERSION="1.1.0"

echo "==> Cleaning up old builds..."
rm -rf ./release
mkdir -p ./release

echo "==> Build for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./release/aes128-cli-linux-amd64 .

echo "==> Build for Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o ./release/aes128-cli-linux-arm64 .

echo "==> Generating checksums..."
cd ./release
sha256sum aes128-cli-linux-amd64 aes128-cli-linux-arm64 > checksums.txt
cd ..

echo "==> Build complete! Files are located in the ./release folder."
