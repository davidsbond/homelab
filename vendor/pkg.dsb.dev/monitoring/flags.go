package monitoring

import "github.com/urfave/cli/v2"

// Flags contains all command-line flags that can be used to configure monitoring.
var Flags = []cli.Flag{
	&cli.BoolFlag{
		Name:        "monitoring-disabled",
		Usage:       "Disables application monitoring",
		EnvVars:     []string{"MONITORING_DISABLED"},
		Destination: &config.disabled,
	},
	&cli.StringFlag{
		Name:        "monitoring-environment",
		Usage:       "Environment to use when writing to application monitoring",
		EnvVars:     []string{"MONITORING_ENVIRONMENT"},
		Destination: &config.environment,
		Value:       "development",
	},
}
