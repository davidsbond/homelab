terraform {
  required_providers {
    grafana = {
      source  = "grafana/grafana"
      version = "1.8.1"
    }
  }
}

provider "grafana" {
  url = "https://grafana.homelab.dsb.dev"
  auth = join(":", [
    var.grafana_user,
    var.grafana_password
  ])
  org_id = 1
}

resource "grafana_data_source" "jaeger" {
  type = "jaeger"
  name = "Jaeger"
  url  = "http://jaeger.monitoring.svc.cluster.local:16686"
}

resource "grafana_data_source" "prometheus" {
  type       = "prometheus"
  name       = "Prometheus"
  url        = "http://prometheus-server.monitoring.svc.cluster.local"
  is_default = true
}

resource "grafana_dashboard" "kubernetes" {
  config_json = file("${path.module}/dashboards/kubernetes.json")
}

resource "grafana_dashboard" "devices" {
  config_json = file("${path.module}/dashboards/devices.json")
}

resource "grafana_dashboard" "application_health" {
  config_json = file("${path.module}/dashboards/application_health.json")
}

resource "grafana_dashboard" "kubernetes_volumes" {
  config_json = file("${path.module}/dashboards/kubernetes_volumes.json")
}

resource "grafana_dashboard" "nas" {
  config_json = file("${path.module}/dashboards/nas.json")
}

resource "grafana_dashboard" "network" {
  config_json = file("${path.module}/dashboards/network.json")
}

resource "grafana_dashboard" "weather" {
  config_json = file("${path.module}/dashboards/weather.json")
}

resource "grafana_dashboard" "minecraft" {
  config_json = file("${path.module}/dashboards/minecraft.json")
}
