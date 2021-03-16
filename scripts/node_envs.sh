#!/usr/bin/env bash

TAILSCALE_ADDR=$(tailscale status -json | jq -r .Self.TailAddr)
K3S_NODE_NAME=k3s-$(echo "$TAILSCALE_ADDR" | tr . -)

export TAILSCALE_ADDR
export K3S_NODE_NAME
