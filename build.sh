#!/bin/bash
set -e

VERSION="1.1.0"

echo "==> Очистка старых сборок..."
rm -rf ./release
mkdir -p ./release

echo "==> Сборка для Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./release/aes128-cli-linux-amd64 .

echo "==> Сборка для Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o ./release/aes128-cli-linux-arm64 .

echo "==> Генерация контрольных сумм..."
cd ./release
sha256sum aes128-cli-linux-amd64 aes128-cli-linux-arm64 > checksums.txt
cd ..

echo "==> Сборка завершена! Файлы находятся в папке ./release"
