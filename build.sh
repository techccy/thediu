#!/bin/bash

VERSION="v1.0.1"
TIME=$(date '+%Y-%m-%d %H:%M:%S')
LDFLAGS="-X 'main.Version=${VERSION}' -X 'main.BuildTime=${TIME}' -s -w"

echo "🔨 开始编译 CCY CLI ${VERSION} ..."

# 清理旧目录
rm -rf build/
mkdir -p build

echo "🍎 编译 macOS..."
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="${LDFLAGS}" -o build/ccy-mac-arm64 ./main.go

echo "🐧 编译 Linux..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o build/ccy-linux-amd64 ./main.go

echo "🪟 编译 Windows..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o build/ccy-win.exe ./main.go

echo "✅ 编译完成！产物已存放在 ./build 目录下。"