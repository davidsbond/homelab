package environment

import (
	"log"
	"strconv"
	"time"
)

var (
	// ApplicationName is the name of the application, injected at build time.
	ApplicationName string

	// ApplicationDescription is the description of the application, injected at build time.
	ApplicationDescription string

	// Version is the version of the application, injected at build time.
	Version string

	// The timestamp the application was compiled, injected at build time.
	compiled string
)

// Compiled returns the time.Time representation of the application's
// build timestamp.
func Compiled() time.Time {
	if compiled == "" {
		return time.Time{}
	}

	unix, err := strconv.ParseInt(compiled, 10, 64)
	if err != nil {
		log.Println(err.Error())
		return time.Time{}
	}

	return time.Unix(unix, 0)
}
