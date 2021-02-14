package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"pkg.dsb.dev/app"
	"pkg.dsb.dev/closers"
	"pkg.dsb.dev/flag"
	"pkg.dsb.dev/logging"
	"pkg.dsb.dev/storage/blob"
	"pkg.dsb.dev/storage/ftp"
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(
			&flag.String{
				Name:        "ftp-addr",
				Usage:       "Address of the FTP server",
				Destination: &ftpAddress,
				EnvVar:      "FTP_ADDRESS",
				Required:    true,
			},
			&flag.String{
				Name:        "ftp-user",
				Usage:       "Username for authenticating with the FTP server",
				Destination: &ftpUser,
				EnvVar:      "FTP_USER",
				Required:    true,
			},
			&flag.String{
				Name:        "ftp-password",
				Usage:       "Password for authenticating with the FTP server",
				Destination: &ftpPassword,
				EnvVar:      "FTP_PASSWORD",
				Required:    true,
			},
			&flag.String{
				Name:        "ftp-path",
				Usage:       "Path to recursively write to the blob store",
				Destination: &ftpPath,
				EnvVar:      "FTP_PATH",
				Value:       "/",
			},
			&flag.String{
				Name:        "bucket-dsn",
				Usage:       "DSN for the bucket to place the backup",
				EnvVar:      "BUCKET_DSN",
				Destination: &bucketDSN,
				Required:    true,
			},
			&flag.Boolean{
				Name:        "ignore-ftp-errors",
				Usage:       "If true, continues iterating over files in the FTP server on error",
				EnvVar:      "IGNORE_FTP_ERRORS",
				Destination: &ignoreFTPErrors,
			},
		),
	)

	if err := a.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var (
	ftpAddress      string
	ftpUser         string
	ftpPassword     string
	ftpPath         string
	bucketDSN       string
	ignoreFTPErrors bool
)

func run(ctx context.Context) error {
	bkt, err := blob.OpenBucket(ctx, bucketDSN)
	if err != nil {
		return err
	}
	defer closers.Close(bkt)

	conn, err := ftp.Open(ctx, ftpAddress, ftp.WithCredentials(ftpUser, ftpPassword))
	if err != nil {
		return err
	}
	defer closers.Close(conn)

	return syncFiles(ctx, bkt, conn)
}

func syncFiles(ctx context.Context, bkt *blob.Bucket, conn *ftp.Conn) error {
	return conn.Walk(ctx, ftpPath, func(path string, info os.FileInfo, err error) error {
		switch {
		case err != nil:
			return err
		case info.IsDir():
			return nil
		}

		wr, err := bkt.NewWriter(ctx, path)
		if err != nil {
			return err
		}
		defer closers.Close(wr)

		rd, err := conn.NewReader(path)
		if err != nil {
			return fmt.Errorf("failed to open reader: %w", err)
		}
		defer closers.Close(rd)

		_, err = io.Copy(wr, rd)
		switch {
		case err != nil && ignoreFTPErrors:
			logging.WithError(err).Error("failed to write file")
			return nil
		case err != nil:
			return fmt.Errorf("failed to write file: %w", err)
		default:
			return nil
		}
	})
}
