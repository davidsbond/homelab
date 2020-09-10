package main

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "network"
	subsystem = "speed"
)

var (
	latency = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "latency",
		Help:      "HTTP request latency",
	})

	upload = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "upload",
		Help:      "HTTP upload speed",
	})

	download = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "download",
		Help:      "HTTP download speed",
	})
)
