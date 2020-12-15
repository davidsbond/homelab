package environment

import "pkg.dsb.dev/flag"

// Flags contains all command-line flags that can be used to configure the application environment.
var Flags = flag.Flags{
	&flag.Boolean{
		Name:        "auto-max-procs-disabled",
		Usage:       "Disables automatically setting GOMAXPROCS to the linux CPU quota",
		EnvVar:      "AUTO_MAX_PROCS_DISABLED",
		Destination: &autoMaxProcsDisabled,
	},
}

var autoMaxProcsDisabled bool
