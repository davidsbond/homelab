package logging

import (
	"pkg.dsb.dev/flag"
)

// Flags contains all command-line flags used by the logger.
var Flags = flag.Flags{
	&flag.String{
		Name:        "log-level",
		Usage:       "Sets the verbosity of logs (panic, fatal, error, warn, info, debug & trace)",
		EnvVar:      "LOG_LEVEL",
		Value:       "error",
		Destination: &level,
	},
	&flag.String{
		Name:        "log-format",
		Usage:       "Sets the desired log formatter (text, json)",
		EnvVar:      "LOG_FORMAT",
		Value:       "text",
		Destination: &format,
	},
}
