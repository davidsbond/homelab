package main

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "world"
	subsystem = "ping"
)

var averageRTT = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: namespace,
	Subsystem: subsystem,
	Name:      "average_rtt",
	Help:      "Average round-trip time",
}, []string{"name", "country_code"})
