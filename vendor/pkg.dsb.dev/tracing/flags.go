package tracing

import (
	"time"

	"github.com/urfave/cli/v2"
)

// Flags contains all command-line flags that can be used to configure tracing.
var Flags = []cli.Flag{
	&cli.BoolFlag{
		Name:        "tracer-disabled",
		Usage:       "If true, disables opentracing",
		EnvVars:     []string{"TRACER_DISABLED"},
		Destination: &config.disabled,
	},
	&cli.StringFlag{
		Name:        "tracer-host",
		Usage:       "Host for opentracing",
		EnvVars:     []string{"TRACER_HOST"},
		Value:       "jaeger:6831",
		Destination: &config.host,
	},
	&cli.Float64Flag{
		Name:        "tracer-sample-rate",
		Usage:       "Sample rate for traces",
		EnvVars:     []string{"TRACER_SAMPLE_RATE"},
		Value:       1,
		Destination: &config.sampleRate,
	},
	&cli.DurationFlag{
		Name:        "tracer-flush-interval",
		Usage:       "Buffer flushing interval for traces",
		EnvVars:     []string{"TRACER_BUFFER_FLUSH_INTERVAL"},
		Value:       time.Second,
		Destination: &config.bufferFlushInterval,
	},
}
