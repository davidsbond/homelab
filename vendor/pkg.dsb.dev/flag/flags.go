// Package flag contains types that represent typed command-line flags.
package flag

import (
	"time"

	"github.com/urfave/cli/v2"
)

type (
	// The Flag interface describes types that can be unwrapped into a cli.Flag
	// implementation.
	Flag interface {
		Unwrap() cli.Flag
	}

	// The Flags type is a collection of Flag implementations.
	Flags []Flag
)

// Unwrap all flags in the collection to their cli.Flag equivalents.
func (f Flags) Unwrap() []cli.Flag {
	out := make([]cli.Flag, len(f))
	for i, flag := range f {
		out[i] = flag.Unwrap()
	}
	return out
}

// The String type represents a command-line flag that is parsed as a string value.
type String struct {
	Name        string
	Usage       string
	Value       string
	Destination *string
	EnvVar      string
	Required    bool
}

// Unwrap the String into its cli.Flag equivalent.
func (f *String) Unwrap() cli.Flag {
	fl := &cli.StringFlag{
		Name:        f.Name,
		Usage:       f.Usage,
		Value:       f.Value,
		Destination: f.Destination,
		Required:    f.Required,
	}

	if f.EnvVar != "" {
		fl.EnvVars = []string{f.EnvVar}
	}

	return fl
}

// The Boolean type represents a command-line flag that is parsed as a boolean value.
type Boolean struct {
	Value       bool
	Required    bool
	Destination *bool
	Name        string
	Usage       string
	EnvVar      string
}

// Unwrap the Boolean into its cli.Flag equivalent.
func (f *Boolean) Unwrap() cli.Flag {
	fl := &cli.BoolFlag{
		Name:        f.Name,
		Usage:       f.Usage,
		Value:       f.Value,
		Destination: f.Destination,
		Required:    f.Required,
	}

	if f.EnvVar != "" {
		fl.EnvVars = []string{f.EnvVar}
	}

	return fl
}

// The Float64 type represents a command-line flag that is parsed as a float64 value.
type Float64 struct {
	Name        string
	Usage       string
	Value       float64
	Destination *float64
	EnvVar      string
	Required    bool
}

// Unwrap the Float64 into its cli.Flag equivalent.
func (f *Float64) Unwrap() cli.Flag {
	fl := &cli.Float64Flag{
		Name:        f.Name,
		Usage:       f.Usage,
		Value:       f.Value,
		Destination: f.Destination,
		Required:    f.Required,
	}

	if f.EnvVar != "" {
		fl.EnvVars = []string{f.EnvVar}
	}

	return fl
}

// The Duration type represents a command-line flag that is parsed as a time.Duration value.
type Duration struct {
	Name        string
	Usage       string
	Value       time.Duration
	Destination *time.Duration
	EnvVar      string
	Required    bool
}

// Unwrap the Duration into its cli.Flag equivalent.
func (f *Duration) Unwrap() cli.Flag {
	fl := &cli.DurationFlag{
		Name:        f.Name,
		Usage:       f.Usage,
		Value:       f.Value,
		Destination: f.Destination,
		Required:    f.Required,
	}

	if f.EnvVar != "" {
		fl.EnvVars = []string{f.EnvVar}
	}

	return fl
}
