package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"pkg.dsb.dev/app"
	"pkg.dsb.dev/cron"
	"pkg.dsb.dev/metrics"

	"github.com/davidsbond/homelab/internal/worldping"
)

const (
	defaultFrequency = time.Hour / 4
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(
			&cli.DurationFlag{
				Name:        "frequency",
				Usage:       "How often to run the ping test.",
				EnvVars:     []string{"FREQUENCY"},
				Value:       defaultFrequency,
				Destination: &frequency,
			},
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

var (
	frequency  time.Duration
	privileged bool
)

func run(ctx context.Context) error {
	metrics.Register(averageRTT)

	return cron.Every(ctx, frequency, func(ctx context.Context) error {
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

		return nil
	})
}
