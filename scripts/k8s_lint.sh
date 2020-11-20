#!/usr/bin/env bash

DIR=$(pwd)
MANIFEST_DIR=${DIR}/manifests
FILES=$(find "${MANIFEST_DIR}" -type d \( -path "**/secrets" -o -path "**/config" \) -prune -false -o -name '*.yaml' | sort | uniq)

kubeval --ignore-missing-schemas --strict $FILES

