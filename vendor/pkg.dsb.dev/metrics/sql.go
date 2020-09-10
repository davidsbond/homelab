package metrics

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	sqlNamespace = "sql"
	subsystem    = "db"
)

type sqlStatsCollector struct {
	db *sql.DB

	// descriptions of exported metrics
	maxOpenDesc           *prometheus.Desc
	openDesc              *prometheus.Desc
	inUseDesc             *prometheus.Desc
	idleDesc              *prometheus.Desc
	waitedForDesc         *prometheus.Desc
	blockedSecondsDesc    *prometheus.Desc
	closedMaxIdleDesc     *prometheus.Desc
	closedMaxLifetimeDesc *prometheus.Desc
}

// newSQLStatsCollector creates a new sqlStatsCollector.
func newSQLStatsCollector(db *sql.DB) *sqlStatsCollector {
	labels := prometheus.Labels{}
	return &sqlStatsCollector{
		db: db,
		maxOpenDesc: prometheus.NewDesc(
			prometheus.BuildFQName(sqlNamespace, subsystem, "max_open"),
			"Maximum number of open connections to the database.",
			nil,
			labels,
		),
		openDesc: prometheus.NewDesc(
			prometheus.BuildFQName(sqlNamespace, subsystem, "open"),
			"The number of established connections both in use and idle.",
			nil,
			labels,
		),
		inUseDesc: prometheus.NewDesc(
			prometheus.BuildFQName(sqlNamespace, subsystem, "in_use"),
			"The number of connections currently in use.",
			nil,
			labels,
		),
		idleDesc: prometheus.NewDesc(
			prometheus.BuildFQName(sqlNamespace, subsystem, "idle"),
			"The number of idle connections.",
			nil,
			labels,
		),
		waitedForDesc: prometheus.NewDesc(
			prometheus.BuildFQName(sqlNamespace, subsystem, "waited_for"),
			"The total number of connections waited for.",
			nil,
			labels,
		),
		blockedSecondsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(sqlNamespace, subsystem, "blocked_seconds"),
			"The total time blocked waiting for a new connection.",
			nil,
			labels,
		),
		closedMaxIdleDesc: prometheus.NewDesc(
			prometheus.BuildFQName(sqlNamespace, subsystem, "closed_max_idle"),
			"The total number of connections closed due to SetMaxIdleConns.",
			nil,
			labels,
		),
		closedMaxLifetimeDesc: prometheus.NewDesc(
			prometheus.BuildFQName(sqlNamespace, subsystem, "closed_max_lifetime"),
			"The total number of connections closed due to SetConnMaxLifetime.",
			nil,
			labels,
		),
	}
}

// Describe implements the prometheus.Collector interface.
func (c sqlStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.maxOpenDesc
	ch <- c.openDesc
	ch <- c.inUseDesc
	ch <- c.idleDesc
	ch <- c.waitedForDesc
	ch <- c.blockedSecondsDesc
	ch <- c.closedMaxIdleDesc
	ch <- c.closedMaxLifetimeDesc
}

// Collect implements the prometheus.Collector interface.
func (c sqlStatsCollector) Collect(ch chan<- prometheus.Metric) {
	stats := c.db.Stats()

	ch <- prometheus.MustNewConstMetric(
		c.maxOpenDesc,
		prometheus.GaugeValue,
		float64(stats.MaxOpenConnections),
	)
	ch <- prometheus.MustNewConstMetric(
		c.openDesc,
		prometheus.GaugeValue,
		float64(stats.OpenConnections),
	)
	ch <- prometheus.MustNewConstMetric(
		c.inUseDesc,
		prometheus.GaugeValue,
		float64(stats.InUse),
	)
	ch <- prometheus.MustNewConstMetric(
		c.idleDesc,
		prometheus.GaugeValue,
		float64(stats.Idle),
	)
	ch <- prometheus.MustNewConstMetric(
		c.waitedForDesc,
		prometheus.CounterValue,
		float64(stats.WaitCount),
	)
	ch <- prometheus.MustNewConstMetric(
		c.blockedSecondsDesc,
		prometheus.CounterValue,
		stats.WaitDuration.Seconds(),
	)
	ch <- prometheus.MustNewConstMetric(
		c.closedMaxIdleDesc,
		prometheus.CounterValue,
		float64(stats.MaxIdleClosed),
	)
	ch <- prometheus.MustNewConstMetric(
		c.closedMaxLifetimeDesc,
		prometheus.CounterValue,
		float64(stats.MaxLifetimeClosed),
	)
}
