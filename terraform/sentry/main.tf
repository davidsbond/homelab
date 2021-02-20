terraform {
  required_providers {
    sentry = {
      source  = "jianyuan/sentry"
      version = "0.6.0"
    }
  }
}

provider "sentry" {
  token = var.sentry_token
}

resource "sentry_team" "default" {
  organization = var.sentry_organization
  name         = "Default"
}

resource "sentry_project" "homelab" {
  organization = var.sentry_organization
  name         = "Homelab"
  team         = sentry_team.default.id
  resolve_age  = 720
  platform     = "go"
}
