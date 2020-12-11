module "cloudflare" {
  source             = "./cloudflare"
  cloudflare_api_key = var.cloudflare_api_key
  cloudflare_email   = var.cloudflare_email
  homelab_ip         = var.homelab_ip
  nas_ip             = var.nas_ip
}

module "minio" {
  source           = "./minio"
  minio_access_key = var.minio_access_key
  minio_secret_key = var.minio_secret_key
  minio_url        = var.minio_url
}
