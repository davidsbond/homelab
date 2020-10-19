package metrics

import (
	"pkg.dsb.dev/flag"
)

// Flags contains all command-line flags that can be used to configure metrics.
var Flags = flag.Flags{
	&flag.Boolean{
		Name:        "metrics-disabled",
		Usage:       "Disables exporting prometheus metrics",
		EnvVar:      "METRICS_DISABLED",
		Destination: &disabled,
	},
	&flag.String{
		Name:        "metrics-push-url",
		Usage:       "URL of the prometheus push gateway, if set, metrics are pushed",
		EnvVar:      "METRICS_PUSH_URL",
		Destination: &pushURL,
	},
}
