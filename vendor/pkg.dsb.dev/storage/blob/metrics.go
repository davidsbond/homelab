package blob

import (
	"github.com/prometheus/client_golang/prometheus"

	"pkg.dsb.dev/metrics"
)

const (
	namespace = "blob"
	subsystem = "bucket"
)

func init() {
	metrics.Register(bytesRead, bytesWritten, blobsOpen)
}

var (
	bytesWritten = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "bytes_written",
		Help:      "Total number of bytes written to blob storage",
	})

	bytesRead = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "bytes_read",
		Help:      "Total number of bytes read from blob storage",
	})

	blobsOpen = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "open_connections",
		Help:      "Total number of open blob readers/writers",
	})
)
