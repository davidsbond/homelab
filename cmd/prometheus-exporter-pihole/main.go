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

	"github.com/davidsbond/homelab/internal/pihole"
)

const (
	defaultFrequency = time.Hour / 2
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
				Name:        "pihole-url",
				Usage:       "The URL of the pihole instance",
				EnvVars:     []string{"PIHOLE_URL"},
				Required:    true,
				Destination: &piHoleURL,
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
	piHoleURL string
)

func run(ctx context.Context) error {
	metrics.Register(
		adsBlocked,
		adsPercentage,
		clientsEverSeen,
		dnsQueries,
		domainsBlocked,
		queriesCached,
		queriesForwarded,
		uniqueClients,
		uniqueDomains,
	)

	ph, err := pihole.New(piHoleURL)
	if err != nil {
		return err
	}

	return cron.Every(ctx, frequency, func(ctx context.Context) error {
		summary, err := ph.Summary(ctx)
		if err != nil {
			return err
		}

		adsBlocked.Set(summary.AdsBlocked)
		adsPercentage.Set(summary.AdsPercentage)
		clientsEverSeen.Set(summary.ClientsEverSeen)
		dnsQueries.Set(summary.DNSQueries)
		domainsBlocked.Set(summary.DomainsBlocked)
		queriesCached.Set(summary.QueriesCached)
		uniqueClients.Set(summary.UniqueClients)
		uniqueDomains.Set(summary.UniqueDomains)
		return nil
	})
}
