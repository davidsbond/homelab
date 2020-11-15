#!/usr/bin/env bash

GO_BIN_DIR=$(go env GOROOT)/bin
GO_ARCH=$(go env GOARCH)

UPX_VERSION=$(git ls-remote --tags --refs --sort="v:refname" https://github.com/upx/upx.git | tail -n1 | sed 's/.*\///')
UPX_URL=https://github.com/upx/upx/releases/download/${UPX_VERSION}/upx-${UPX_VERSION:1}-${GO_ARCH}_linux.tar.xz

# Download latest UPX release and unarchive
wget "${UPX_URL}" -O upx.tar.xz
tar -xf upx.tar.xz

# Move UPX binary into go bin directory
mv upx-"${UPX_VERSION:1}"-"${GO_ARCH}"_linux/upx "${GO_BIN_DIR}"

# Tidy up after ourselves
rm upx.tar.xz
rm -rf upx-"${UPX_VERSION:1}"-"${GO_ARCH}"_linux
