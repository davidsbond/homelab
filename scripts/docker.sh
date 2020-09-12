#!/usr/bin/env bash

TAG=$(git describe --tags --always)

docker build \
  -t docker.pkg.github.com/davidsbond/homelab/homelab:latest \
  -t docker.pkg.github.com/davidsbond/homelab/homelab:${TAG} .
