package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
	"pkg.dsb.dev/app"
	"pkg.dsb.dev/closers"
	"pkg.dsb.dev/flag"
	"pkg.dsb.dev/storage/blob"

	"github.com/davidsbond/homelab/internal/grafana"
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(
			&flag.String{
				Name:        "grafana-url",
				Usage:       "URL of the grafana instance to backup",
				EnvVar:      "GRAFANA_URL",
				Destination: &grafanaURL,
				Required:    true,
			},
			&flag.String{
				Name:        "grafana-api-key",
				Usage:       "API key to use for authenticating with grafana",
				EnvVar:      "GRAFANA_API_KEY",
				Destination: &grafanaAPIKey,
				Required:    true,
			},
			&flag.String{
				Name:        "bucket-dsn",
				Usage:       "DSN for the bucket to place the backup",
				EnvVar:      "BUCKET_DSN",
				Destination: &bucketDSN,
				Required:    true,
			},
		),
	)

	if err := a.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var (
	grafanaURL    string
	grafanaAPIKey string
	bucketDSN     string
)

func run(ctx context.Context) error {
	bkt, err := blob.OpenBucket(ctx, bucketDSN)
	if err != nil {
		return err
	}
	defer closers.Close(bkt)

	client, err := grafana.NewClient(grafanaURL, grafanaAPIKey)
	if err != nil {
		return err
	}

	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return backupDashboards(ctx, client, bkt)
	})
	grp.Go(func() error {
		return backupDataSources(ctx, client, bkt)
	})

	return grp.Wait()
}

const (
	ext           = ".json"
	dashboardDir  = "dashboard"
	dataSourceDir = "datasource"
)

func backupDashboards(ctx context.Context, cl *grafana.Client, bkt *blob.Bucket) error {
	return cl.IterateDashboards(ctx, func(ctx context.Context, d *grafana.Dashboard) error {
		key := filepath.Join(dashboardDir, d.UID()+ext)
		wr, err := bkt.NewWriter(ctx, key)
		if err != nil {
			return err
		}
		defer closers.Close(wr)

		_, err = io.Copy(wr, d)
		return err
	})
}

func backupDataSources(ctx context.Context, cl *grafana.Client, bkt *blob.Bucket) error {
	return cl.IterateDataSources(ctx, func(ctx context.Context, d *grafana.DataSource) error {
		key := filepath.Join(dataSourceDir, d.ID()+ext)
		wr, err := bkt.NewWriter(ctx, key)
		if err != nil {
			return err
		}
		defer closers.Close(wr)

		_, err = io.Copy(wr, d)
		return err
	})
}
