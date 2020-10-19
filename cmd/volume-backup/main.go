package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"pkg.dsb.dev/app"
	"pkg.dsb.dev/closers"
	"pkg.dsb.dev/cron"
	"pkg.dsb.dev/flag"
	"pkg.dsb.dev/storage/blob"

	"github.com/davidsbond/homelab/internal/filesystem"
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(
			&flag.String{
				Name:        "volume-dir",
				Usage:       "Directory to back up",
				EnvVar:      "VOLUME_DIR",
				Destination: &volumeDir,
				Required:    true,
			},
			&flag.String{
				Name:        "bucket-dsn",
				Usage:       "DSN for the bucket to place the backup",
				EnvVar:      "BUCKET_DSN",
				Destination: &bucketDSN,
				Required:    true,
			},
			&flag.String{
				Name:        "bucket-dir",
				Usage:       "Location in the bucket to place the backup",
				EnvVar:      "BUCKET_DIR",
				Destination: &bucketDir,
				Required:    true,
			},
			&flag.Duration{
				Name:        "frequency",
				Usage:       "How often to perform a backup",
				EnvVar:      "FREQUENCY",
				Destination: &frequency,
				Value:       time.Hour,
			},
		),
	)

	if err := a.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var (
	volumeDir string
	bucketDSN string
	bucketDir string
	frequency time.Duration
)

func run(ctx context.Context) error {
	bkt, err := blob.OpenBucket(ctx, bucketDSN)
	if err != nil {
		return err
	}
	defer closers.Close(bkt)

	return cron.Every(ctx, frequency, func(ctx context.Context) error {
		key := filepath.Join(bucketDir, time.Now().Format("2006-01-02.tar.gz"))
		wr, err := bkt.NewWriter(ctx, key)
		if err != nil {
			return err
		}
		defer closers.Close(wr)

		return filesystem.Archive(ctx, volumeDir, wr)
	})
}
