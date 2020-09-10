package monitoring

import "time"

type (
	options struct {
		dsn          string
		flushTimeout time.Duration
		environment  string
		disabled     bool
	}
)

var config = &options{
	dsn:          "https://04a1f22efdb24c659631e5769b78a8fc@o407749.ingest.sentry.io/5277353",
	flushTimeout: time.Minute,
	environment:  "development",
}
