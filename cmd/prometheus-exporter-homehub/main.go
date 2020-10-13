package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"pkg.dsb.dev/app"
	"pkg.dsb.dev/metrics"

	"github.com/davidsbond/homelab/internal/homehub"
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(
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

var homeHubURL string

func run(ctx context.Context) error {
	metrics.Register(uptime, bytesUp, bytesDown)

	hh, err := homehub.New(homeHubURL)
	if err != nil {
		return err
	}

	results, err := hh.Summary(ctx)
	if err != nil {
		return err
	}

	uptime.Set(results.Uptime)
	bytesDown.Set(results.BytesDown)
	bytesUp.Set(results.BytesUp)

	return metrics.Push()
}
