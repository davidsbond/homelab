package main

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "synology"
	subsystem = "stats"
)

var (
	diskSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "disk_size",
		Help:      "The size of the disk in bytes.",
	}, []string{"name"})

	diskTemp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "disk_temperature",
		Help:      "The temperature of the disk.",
	}, []string{"name"})

	volumeSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "volume_size",
		Help:      "The size of the volume in bytes.",
	}, []string{"name"})

	volumeUsed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "volume_bytes_used",
		Help:      "The amount of space used on the volume, in bytes.",
	}, []string{"name"})

	uptime = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "uptime",
		Help:      "Uptime of the nas in seconds",
	})
)
