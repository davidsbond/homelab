package health

import "github.com/urfave/cli/v2"

// Flags contains all command-line flags that can be used to configure health checks.
var Flags = []cli.Flag{
	&cli.BoolFlag{
		Name:        "health-check-disabled",
		Usage:       "Disables the health check endpoint",
		EnvVars:     []string{"HEALTH_CHECK_DISABLED"},
		Destination: &disabled,
	},
}
