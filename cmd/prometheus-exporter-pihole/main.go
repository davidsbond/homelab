package main

import (
	"context"
	"fmt"
	"os"

	"pkg.dsb.dev/app"
	"pkg.dsb.dev/flag"
	"pkg.dsb.dev/metrics"

	"github.com/davidsbond/homelab/internal/pihole"
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(
			&flag.String{
				Name:        "pihole-url",
				Usage:       "The URL of the pihole instance",
				EnvVar:      "PIHOLE_URL",
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

var piHoleURL string

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

	return metrics.Push()
}
