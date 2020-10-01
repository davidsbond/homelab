#!/usr/bin/env bash

TAG=$(git describe --tags --always)

docker build \
  -t docker.pkg.github.com/davidsbond/homelab/homelab:latest \
  -t docker.pkg.github.com/davidsbond/homelab/homelab:${TAG} .

# If env variable is set, auth with registry and push images.
if [ "$PUSH" = "true" ]; then
  echo ${GITHUB_TOKEN} | docker login docker.pkg.github.com -u davidsbond --password-stdin
  docker push docker.pkg.github.com/davidsbond/homelab/homelab:latest
  docker push docker.pkg.github.com/davidsbond/homelab/homelab:${TAG}
fi
