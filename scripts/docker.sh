#!/usr/bin/env bash

REGISTRY=ghcr.io
IMAGE=${REGISTRY}/davidsbond/homelab
TAG=$(git describe --tags --always)

docker build \
  -t ${IMAGE}:latest \
  -t ${IMAGE}:${TAG} .

# If env variable is set, auth with registry and push images.
if [ "$PUSH" = "true" ]; then
  echo ${GITHUB_TOKEN} | docker login ${REGISTRY} -u davidsbond --password-stdin
  docker push ${IMAGE}:latest
  docker push ${IMAGE}:${TAG}
fi
