package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"pkg.dsb.dev/app"
	"pkg.dsb.dev/metrics"

	"github.com/davidsbond/homelab/internal/worldping"
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(
			&cli.BoolFlag{
				Name:        "privileged",
				Usage:       "If true, uses privileged ICMP requests.",
				EnvVars:     []string{"PRIVILEGED"},
				Destination: &privileged,
			},
		),
	)

	if err := a.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var privileged bool

func run(ctx context.Context) error {
	metrics.Register(averageRTT)

	results, err := worldping.Run(ctx, privileged)
	if err != nil {
		return err
	}

	for server, result := range results {
		// Zero result indicates a timeout or error.
		if result == 0 {
			continue
		}

		averageRTT.
			WithLabelValues(server.Name, server.Code).
			Set(float64(result.Milliseconds()))
	}

	return metrics.Push()
}
