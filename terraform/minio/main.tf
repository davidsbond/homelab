terraform {
  required_providers {
    minio = {
      source  = "registry.terraform.io/aminueza/minio"
      version = "1.2.0"
    }
  }
}

provider "minio" {
  minio_server     = var.minio_url
  minio_region     = "us-east-1"
  minio_access_key = var.minio_access_key
  minio_secret_key = var.minio_secret_key
}

resource "minio_s3_bucket" "database_backups" {
  bucket = "databases"
  acl    = "public"
}

resource "minio_s3_bucket" "grafana_backups" {
  bucket = "grafana"
  acl    = "public"
}

resource "minio_s3_bucket" "bytecrypt_data" {
  bucket = "bytecrypt"
  acl    = "public"
}

resource "minio_s3_bucket" "minecraft_data" {
  bucket = "minecraft"
  acl    = "public"
}

