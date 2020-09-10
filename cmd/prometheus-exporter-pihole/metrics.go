package main

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "pihole"
	subsystem = "summary"
)

var (
	adsBlocked = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "ads_blocked",
		Help:      "Number of advertisements blocked today",
	})

	adsPercentage = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "ads_percentage",
		Help:      "Percentage of queries that are blocked adverts",
	})

	clientsEverSeen = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "clients_ever_seen",
		Help:      "Count of clients ever seen by pihole",
	})

	dnsQueries = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "dns_queries_today",
		Help:      "Count of DNS queries made today",
	})

	domainsBlocked = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "domains_blocked",
		Help:      "Count of blocked domains",
	})

	queriesCached = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "queries_cached",
		Help:      "Count of queries cached",
	})

	queriesForwarded = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "queries_forwarded",
		Help:      "Count of queries forwarded",
	})

	uniqueClients = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "unique_clients",
		Help:      "Count of unique clients",
	})

	uniqueDomains = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "unique_domains",
		Help:      "Count of unique domains",
	})
)
