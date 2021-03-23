package tracing

import (
	"time"

	"pkg.dsb.dev/flag"
)

// Flags contains all command-line flags that can be used to configure tracing.
var Flags = flag.Flags{
	&flag.Boolean{
		Name:        "tracer-disabled",
		Usage:       "If true, disables opentracing",
		EnvVar:      "TRACER_DISABLED",
		Destination: &config.disabled,
		Hidden:      true,
	},
	&flag.String{
		Name:        "tracer-host",
		Usage:       "Host for opentracing",
		EnvVar:      "TRACER_HOST",
		Destination: &config.host,
		Hidden:      true,
	},
	&flag.Float64{
		Name:        "tracer-sample-rate",
		Usage:       "Sample rate for traces",
		EnvVar:      "TRACER_SAMPLE_RATE",
		Value:       1,
		Destination: &config.sampleRate,
		Hidden:      true,
	},
	&flag.Duration{
		Name:        "tracer-flush-interval",
		Usage:       "Buffer flushing interval for traces",
		EnvVar:      "TRACER_BUFFER_FLUSH_INTERVAL",
		Value:       time.Second,
		Destination: &config.bufferFlushInterval,
		Hidden:      true,
	},
}
