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
	flushTimeout: time.Minute,
	environment:  "development",
}
