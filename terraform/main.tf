module "cloudflare" {
  source             = "./cloudflare"
  cloudflare_api_key = var.cloudflare_api_key
  cloudflare_email   = var.cloudflare_email
  homelab_0_ip       = var.homelab_0_ip
  homelab_1_ip       = var.homelab_1_ip
  homelab_2_ip       = var.homelab_2_ip
  homelab_3_ip       = var.homelab_3_ip
  nas_ip             = var.nas_ip
}

module "minio" {
  source           = "./minio"
  minio_access_key = var.minio_access_key
  minio_secret_key = var.minio_secret_key
  minio_url        = var.minio_url
}

module "grafana" {
  source           = "./grafana"
  grafana_password = var.grafana_password
  grafana_user     = var.grafana_user
}

module "tailscale" {
  source            = "./tailscale"
  tailscale_api_key = var.tailscale_api_key
  tailscale_domain  = var.tailscale_domain
}

module "sentry" {
  source              = "./sentry"
  sentry_token        = var.sentry_token
  sentry_organization = var.sentry_organization
}
