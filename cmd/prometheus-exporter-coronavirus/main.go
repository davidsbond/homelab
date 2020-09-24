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

	"github.com/davidsbond/homelab/internal/coronavirus"
)

const (
	defaultFrequency = time.Hour / 2
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(&cli.DurationFlag{
			Name:        "frequency",
			Usage:       "How often to query coronavirus data.",
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
	metrics.Register(dailyCases, dailyDeaths)

	client, err := coronavirus.New()
	if err != nil {
		return err
	}

	return cron.Every(ctx, frequency, func(ctx context.Context) error {
		results, err := client.GetSummary(ctx, time.Now().AddDate(0, 0, -1))
		if err != nil {
			return err
		}

		dailyCases.Set(results.Cases)
		dailyDeaths.Set(results.Deaths)
		return nil
	})
}
