export TAILSCALE_ADDR=$(tailscale status -json | jq -r .Self.TailAddr)
export K3S_NODE_NAME=k3s-$(echo $TAILSCALE_ADDR | tr . -)
