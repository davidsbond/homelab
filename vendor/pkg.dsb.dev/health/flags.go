package health

import (
	"pkg.dsb.dev/flag"
)

// Flags contains all command-line flags that can be used to configure health checks.
var Flags = flag.Flags{
	&flag.Boolean{
		Name:        "health-check-disabled",
		Usage:       "Disables the health check endpoint",
		EnvVar:      "HEALTH_CHECK_DISABLED",
		Destination: &disabled,
		Hidden:      true,
	},
}
