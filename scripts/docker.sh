#!/usr/bin/env bash

# This script is used to build the docker image that contains all bespoke applications ran
# in the k3s cluster. If no go files have changed, nothing is done here.

CHANGES=$(git diff --name-only HEAD HEAD~1 | grep .go)
if [ -z "$CHANGES" ]
then
  echo "No go files have changed, skipping docker build"
  exit
fi

REGISTRY=ghcr.io
IMAGE=${REGISTRY}/davidsbond/homelab
TAG=$(git describe --tags --always)

echo "${GITHUB_TOKEN}" | docker login ${REGISTRY} -u davidsbond --password-stdin

docker buildx create --use
docker buildx build --platform linux/amd64,linux/arm64 \
  --push \
  -t ${IMAGE}:latest \
  -t ${IMAGE}:"${TAG}" .
