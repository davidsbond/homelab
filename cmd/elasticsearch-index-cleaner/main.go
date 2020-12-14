package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"
	"pkg.dsb.dev/app"
	"pkg.dsb.dev/flag"
	"pkg.dsb.dev/period"

	"github.com/davidsbond/homelab/internal/elasticsearch"
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(
			&flag.String{
				Name:        "elasticsearch-url",
				Usage:       "Elasticsearch host URL",
				EnvVar:      "ELASTICSEARCH_HOST",
				Required:    true,
				Destination: &host,
			},
			&flag.String{
				Name:        "index-formats",
				Usage:       "Space separated index formats, each will be generated with the year, month and day",
				EnvVar:      "INDEX_FORMATS",
				Required:    true,
				Destination: &formats,
			},
			&flag.Duration{
				Name:        "max-age",
				Usage:       "Max age for indices until they're deleted",
				Value:       (period.Day * 3).Duration(),
				Destination: &maxAge,
				EnvVar:      "MAX_AGE",
			},
		),
	)

	if err := a.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var (
	host    string
	formats string
	maxAge  time.Duration
)

func run(ctx context.Context) error {
	client, err := elastic.NewSimpleClient(elastic.SetURL(host))
	if err != nil {
		return err
	}

	return elasticsearch.NewIndexCleaner(client).Clean(ctx, strings.Fields(formats), maxAge)
}
