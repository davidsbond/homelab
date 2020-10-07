#!/usr/bin/env bash

REGISTRY=ghcr.io
IMAGE=${REGISTRY}/davidsbond/homelab
TAG=$(git describe --tags --always)

echo ${GITHUB_TOKEN} | docker login ${REGISTRY} -u davidsbond --password-stdin

docker buildx create --use
docker buildx build --platform linux/amd64,linux/arm64 \
  --push \
  -t ${IMAGE}:latest \
  -t ${IMAGE}:${TAG} .
