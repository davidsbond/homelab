package main

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "filesystem"
	subsystem = "storage"
)

var (
	percentageAvailable = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "percentage_bytes_available",
		Help:      "Available storage space expressed as a percentage.",
	}, []string{"drive"})

	percentageUsed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "percentage_bytes_used",
		Help:      "Used storage space expressed as a percentage.",
	}, []string{"drive"})

	total = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "bytes_total",
		Help:      "Total storage size expressed in bytes.",
	}, []string{"drive"})

	available = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "bytes_available",
		Help:      "Available storage space expressed in bytes.",
	}, []string{"drive"})

	used = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "bytes_used",
		Help:      "Used storage space expressed in bytes.",
	}, []string{"drive"})
)
