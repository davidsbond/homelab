#!/usr/bin/env bash

DIR=$(pwd)
BIN_DIR=${DIR}/bin

# Compress all generated binaries with UPX (https://upx.github.io/)
BINARIES=$(find ${BIN_DIR} -type f | xargs -0 -n1 | sort | uniq)
upx -q ${BINARIES}
