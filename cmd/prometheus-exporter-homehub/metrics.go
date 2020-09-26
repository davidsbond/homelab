package main

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "homehub"
	subsystem = "stats"
)

var (
	uptime = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "uptime",
		Help:      "Uptime of the homehub in seconds",
	})

	bytesUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "bytes_up",
		Help:      "Total number of bytes uploaded",
	})

	bytesDown = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "bytes_down",
		Help:      "Total number of bytes downloaded",
	})
)
