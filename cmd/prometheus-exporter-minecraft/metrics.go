package main

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "minecraft_server"
	subsystem = "stats"
)

var (
	playersOnline = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "players_online",
		Help:      "Number of online players",
	})

	maxPlayers = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "max_players",
		Help:      "Number of maximum online players",
	})
)
