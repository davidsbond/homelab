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
	defaultDrive     = "/"
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
				Value:       defaultDrive,
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

		percentageAvailable.WithLabelValues(drive).Set(results.PercentageAvailable)
		percentageUsed.WithLabelValues(drive).Set(results.PercentageUsed)
		total.WithLabelValues(drive).Set(results.Total)
		available.WithLabelValues(drive).Set(results.Available)
		used.WithLabelValues(drive).Set(results.Used)
		return nil
	})
}
