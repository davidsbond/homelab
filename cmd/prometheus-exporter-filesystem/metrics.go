package main

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "filesystem"
	subsystem = "storage"
)

var (
	percentageAvailable = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "percentage_bytes_available",
		Help:      "Available storage space expressed as a percentage.",
	})

	percentageUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "percentage_bytes_used",
		Help:      "Used storage space expressed as a percentage.",
	})

	total = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "bytes_total",
		Help:      "Total storage size expressed in bytes.",
	})

	available = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "bytes_available",
		Help:      "Available storage space expressed in bytes.",
	})

	used = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "bytes_used",
		Help:      "Used storage space expressed in bytes.",
	})
)
