terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "2.17.0"
    }
  }
}

provider "cloudflare" {
  email   = var.cloudflare_email
  api_key = var.cloudflare_api_key
}

data "cloudflare_zones" "dsb_dev" {
  filter {}
}

resource "cloudflare_record" "homelab_0" {
  zone_id = lookup(data.cloudflare_zones.dsb_dev.zones[0], "id")
  name    = "*.homelab"
  value   = var.homelab_0_ip
  type    = "A"
  ttl     = 3600
}

resource "cloudflare_record" "homelab_1" {
  zone_id = lookup(data.cloudflare_zones.dsb_dev.zones[0], "id")
  name    = "*.homelab"
  value   = var.homelab_1_ip
  type    = "A"
  ttl     = 3600
}

resource "cloudflare_record" "homelab_2" {
  zone_id = lookup(data.cloudflare_zones.dsb_dev.zones[0], "id")
  name    = "*.homelab"
  value   = var.homelab_2_ip
  type    = "A"
  ttl     = 3600
}

resource "cloudflare_record" "homelab_3" {
  zone_id = lookup(data.cloudflare_zones.dsb_dev.zones[0], "id")
  name    = "*.homelab"
  value   = var.homelab_3_ip
  type    = "A"
  ttl     = 3600
}

resource "cloudflare_record" "nas" {
  zone_id = lookup(data.cloudflare_zones.dsb_dev.zones[0], "id")
  name    = "nas"
  value   = var.nas_ip
  type    = "A"
  ttl     = 3600
}
