package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/davidsbond/homelab/internal/homehub"

	"github.com/urfave/cli/v2"
	"pkg.dsb.dev/app"
	"pkg.dsb.dev/cron"
	"pkg.dsb.dev/metrics"
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
				Usage:       "How often to scrape the homehub",
				EnvVars:     []string{"FREQUENCY"},
				Value:       defaultFrequency,
				Destination: &frequency,
			},
			&cli.StringFlag{
				Name:        "homehub-url",
				Usage:       "The URL of the BT homehub",
				EnvVars:     []string{"HOMEHUB_URL"},
				Required:    true,
				Destination: &homeHubURL,
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
	homeHubURL string
)

func run(ctx context.Context) error {
	metrics.Register(uptime, bytesUp, bytesDown)
	hh, err := homehub.New(homeHubURL)
	if err != nil {
		return err
	}

	return cron.Every(ctx, frequency, func(ctx context.Context) error {
		results, err := hh.Summary(ctx)
		if err != nil {
			return err
		}

		uptime.Set(results.Uptime)
		bytesDown.Set(results.BytesDown)
		bytesUp.Set(results.BytesUp)
		return nil
	})
}
