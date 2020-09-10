package metrics

import "github.com/urfave/cli/v2"

// Flags contains all command-line flags that can be used to configure metrics.
var Flags = []cli.Flag{
	&cli.BoolFlag{
		Name:        "metrics-disabled",
		Usage:       "Disables exporting prometheus metrics",
		EnvVars:     []string{"METRICS_DISABLED"},
		Destination: &disabled,
	},
}
