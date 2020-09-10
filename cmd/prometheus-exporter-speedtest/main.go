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

	"github.com/davidsbond/homelab/internal/speedtest"
)

const (
	defaultFrequency = time.Hour / 2
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(&cli.DurationFlag{
			Name:        "frequency",
			Usage:       "How often to run the speed test",
			EnvVars:     []string{"FREQUENCY"},
			Value:       defaultFrequency,
			Destination: &frequency,
		}),
	)

	if err := a.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var frequency time.Duration

func run(ctx context.Context) error {
	metrics.Register(latency, upload, download)

	return cron.Every(ctx, frequency, func(ctx context.Context) error {
		results, err := speedtest.New().Test(ctx)
		if err != nil {
			return err
		}

		latency.Set(results.Latency)
		upload.Set(results.Upload)
		download.Set(results.Download)
		return nil
	})
}
