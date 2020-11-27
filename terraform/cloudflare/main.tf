provider "cloudflare" {
  version = "~> 2.0"
  email   = var.cloudflare_email
  api_key = var.cloudflare_api_key
}

data "cloudflare_zones" "dsb_dev" {
  filter {}
}

resource "cloudflare_record" "homelab" {
  zone_id = lookup(data.cloudflare_zones.dsb_dev.zones[0], "id")
  name    = "*.homelab"
  value   = var.homelab_ip
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
