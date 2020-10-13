package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"pkg.dsb.dev/app"
	"pkg.dsb.dev/metrics"

	"github.com/davidsbond/homelab/internal/coronavirus"
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
	metrics.Register(dailyCases, dailyDeaths)

	client, err := coronavirus.New()
	if err != nil {
		return err
	}

	results, err := client.GetSummary(ctx, time.Now().AddDate(0, 0, -1))
	if err != nil {
		return err
	}

	dailyCases.Set(results.Cases)
	dailyDeaths.Set(results.Deaths)

	return metrics.Push()
}
