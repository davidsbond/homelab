package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/davidsbond/homelab/internal/filesystem"

	"github.com/urfave/cli/v2"
	"pkg.dsb.dev/app"
	"pkg.dsb.dev/cron"
	"pkg.dsb.dev/metrics"
)

const (
	defaultFrequency = time.Minute * 5
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(
			&cli.DurationFlag{
				Name:        "frequency",
				Usage:       "How often to run the speed test",
				EnvVars:     []string{"FREQUENCY"},
				Value:       defaultFrequency,
				Destination: &frequency,
			},
			&cli.StringFlag{
				Name:        "drive",
				Usage:       "Drive to return statistics for",
				EnvVars:     []string{"DRIVE"},
				Destination: &drive,
				Required:    true,
			},
		),
	)

	if err := a.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var (
	frequency time.Duration
	drive     string
)

func run(ctx context.Context) error {
	metrics.Register(
		percentageAvailable,
		percentageUsed,
		total,
		available,
		used,
	)

	return cron.Every(ctx, frequency, func(ctx context.Context) error {
		results, err := filesystem.GetSummary(drive)
		if err != nil {
			return err
		}

		percentageAvailable.Set(results.PercentageAvailable)
		percentageUsed.Set(results.PercentageUsed)
		total.Set(results.Total)
		available.Set(results.Available)
		used.Set(results.Used)
		return nil
	})
}
