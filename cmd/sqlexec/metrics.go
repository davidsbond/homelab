package main

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "sql"
	subsystem = "exec"
)

var (
	rowsAffected = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "rows_affected",
		Help:      "Count of affected rows",
	})

	lastInsertID = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "last_insert_id",
		Help:      "Integer generated by the database in response to a command",
	})
)