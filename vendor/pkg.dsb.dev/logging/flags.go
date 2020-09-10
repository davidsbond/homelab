package logging

import "github.com/urfave/cli/v2"

// Flags contains all command-line flags used by the logger.
var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "log-level",
		Aliases:     nil,
		Usage:       "Sets the verbosity of logs (panic, fatal, error, warn, info, debug & trace)",
		EnvVars:     []string{"LOG_LEVEL"},
		Value:       "error",
		Destination: &level,
	},
}
