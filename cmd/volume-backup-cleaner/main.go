package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"pkg.dsb.dev/app"
	"pkg.dsb.dev/closers"
	"pkg.dsb.dev/flag"
	"pkg.dsb.dev/period"
	"pkg.dsb.dev/storage/blob"
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(
			&flag.String{
				Name:        "bucket-dsn",
				Usage:       "DSN for the bucket to check backups",
				EnvVar:      "BUCKET_DSN",
				Destination: &bucketDSN,
				Required:    true,
			},
			&flag.Duration{
				Name:        "older-than",
				Usage:       "Maximum age of a backup",
				EnvVar:      "OLDER_THAN",
				Destination: &olderThan,
				Value:       period.Week.Duration(),
			},
		),
	)

	if err := a.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var (
	bucketDSN string
	olderThan time.Duration
)

func run(ctx context.Context) error {
	bkt, err := blob.OpenBucket(ctx, bucketDSN)
	if err != nil {
		return err
	}
	defer closers.Close(bkt)

	now := time.Now()
	return bkt.Iterate(ctx, func(ctx context.Context, item blob.Blob) error {
		// Don't do anything if item is less than the configured age.
		if item.ModTime.Add(olderThan).After(now) {
			return nil
		}

		return bkt.Delete(ctx, item.Key)
	})
}
