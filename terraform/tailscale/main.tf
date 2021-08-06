terraform {
  required_providers {
    tailscale = {
      source  = "davidsbond/tailscale"
      version = "0.1.1"
    }
  }
}

provider "tailscale" {
  api_key = var.tailscale_api_key
  tailnet = var.tailscale_domain
}

resource "tailscale_acl" "homelab_acl" {
  acl = file("${path.module}/acl/acl.json")
}
