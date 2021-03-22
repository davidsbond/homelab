package monitoring

import (
	"pkg.dsb.dev/flag"
)

// Flags contains all command-line flags that can be used to configure monitoring.
var Flags = flag.Flags{
	&flag.Boolean{
		Name:        "monitoring-disabled",
		Usage:       "Disables application monitoring",
		EnvVar:      "MONITORING_DISABLED",
		Destination: &config.disabled,
		Hidden:      true,
	},
	&flag.String{
		Name:        "monitoring-environment",
		Usage:       "Environment to use when writing to application monitoring",
		EnvVar:      "MONITORING_ENVIRONMENT",
		Destination: &config.environment,
		Value:       "development",
		Hidden:      true,
	},
	&flag.String{
		Name:        "monitoring-dsn",
		Usage:       "DSN to use for sending reports to sentry",
		EnvVar:      "MONITORING_DSN",
		Destination: &config.dsn,
		Hidden:      true,
	},
}
