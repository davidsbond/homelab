package main

import (
	"context"
	"fmt"
	"os"

	"pkg.dsb.dev/app"
	"pkg.dsb.dev/metrics"

	"github.com/davidsbond/homelab/internal/speedtest"
)

func main() {
	a := app.New(
		app.WithFlags(),
		app.WithRunner(run),
	)

	if err := a.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	metrics.Register(latency, upload, download)

	results, err := speedtest.New().Test(ctx)
	if err != nil {
		return err
	}

	latency.Set(results.Latency)
	upload.Set(results.Upload)
	download.Set(results.Download)

	return metrics.Push()
}
