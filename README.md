# homelab [![PkgGoDev](https://pkg.go.dev/badge/github.com/davidsbond/homelab)](https://pkg.go.dev/github.com/davidsbond/homelab) [![Go Report Card](https://goreportcard.com/badge/github.com/davidsbond/homelab)](https://goreportcard.com/report/github.com/davidsbond/homelab)

Monorepo for my personal homelab. It contains applications and kubernetes manifests for deployment.

## Getting started

This assumes you have [go](https://golang.org/), [docker](https://www.docker.com/) & [make](https://www.gnu.org/software/make/manual/make.html) 
installed along with docker's [buildx](https://docs.docker.com/buildx/working-with-buildx/) plugin.

* Clone the repository
* Install golang tools using `make install-tools`
* Run `make` to build all binaries

## Project structure

* `cmd` - Entry points to any bespoke applications
* `internal` - Packages used throughout the application code
* `manifests` - Kubernetes manifests to run all my homelab applications
* `vendor` - Vendored third-party code

## Third party applications

Here's a list of third-party applications I'm using alongside my custom applications:

* [deluge](https://deluge-torrent.org/) - For torrents
* [home-assistant](https://www.home-assistant.io/) - For interfacing with my ring cameras
* [pihole](https://pi-hole.net/) - For DNS and ad-blocking
* [traefik](https://traefik.io/) - For exposing services on the homelab, usually just their web-interfaces
* [prometheus](https://prometheus.io/) - For scraping all my metrics
* [grafana](https://grafana.com/) - For visualising all my metrics
* [jaeger](https://www.jaegertracing.io/) - For tracing my applications
* [node-exporter](https://github.com/prometheus/node_exporter) - For monitoring the host environment and exporting stats as prometheus metrics
* [helm](https://helm.sh/) - For installing non-trivial k8s applications, like the prometheus operator.

## Prometheus exporters

I've implemented several custom prometheus exporters in this repo that power my dashboards, these are:

* `coronavirus` - Exports UK coronavirus stats as prometheus metrics
* `homehub` - Exports statistics from my BT HomeHub as prometheus metrics
* `pihole` - Exports statistics from my pihole as prometheus metrics
* `speedtest` - Exports [speedtest](https://speedtest.net) results as prometheus metrics
* `weather` - Exports current weather data as prometheus metrics
* `worldping` - Exports world ping times for the local host as prometheus metrics
* `home-assistant` - Proxies prometheus metrics from a home-assistant server.
