module "cloudflare" {
  source             = "./cloudflare"
  cloudflare_api_key = var.cloudflare_api_key
  cloudflare_email   = var.cloudflare_email
  homelab_ip         = var.homelab_ip
  nas_ip             = var.nas_ip
}
