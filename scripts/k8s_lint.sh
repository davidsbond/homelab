#!/usr/bin/env bash

# This script runs validation for all kubernetes manifests, excluding secrets and configuration files
# that are managed with kustomize.

DIR=$(pwd)
MANIFEST_DIR=${DIR}/manifests
FILES=$(find "${MANIFEST_DIR}" -type d \( -path "**/secrets" -o -path "**/config" \) -prune -false -o -name '*.yaml' | sort | uniq)

kubeval --ignore-missing-schemas --strict $FILES

