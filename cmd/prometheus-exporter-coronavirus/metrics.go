package main

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "coronavirus"
	subsystem = "stats"
)

var (
	dailyCases = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "case_count",
		Help:      "Count of daily cases",
	})

	dailyDeaths = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "death_count",
		Help:      "Count of daily deaths",
	})
)
