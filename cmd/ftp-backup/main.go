package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"pkg.dsb.dev/app"
	"pkg.dsb.dev/closers"
	"pkg.dsb.dev/flag"
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
			&flag.String{
				Name:        "zip-name-layout",
				Usage:       "Layout string for the bucket name, should use Go's date strings",
				EnvVar:      "ZIP_NAME_LAYOUT",
				Value:       "2006-01-02.zip",
				Destination: &zipNameLayout,
			},
		),
	)

	if err := a.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var (
	ftpAddress    string
	ftpUser       string
	ftpPassword   string
	ftpPath       string
	bucketDSN     string
	zipNameLayout string
)

func run(ctx context.Context) error {
	bucket, err := blob.OpenBucket(ctx, bucketDSN)
	if err != nil {
		return err
	}
	defer closers.Close(bucket)

	conn, err := ftp.Open(ctx, ftpAddress, ftp.WithCredentials(ftpUser, ftpPassword))
	if err != nil {
		return err
	}
	defer closers.Close(conn)

	key := time.Now().Format(zipNameLayout)
	blobWriter, err := bucket.NewWriter(ctx, key)
	if err != nil {
		return err
	}
	defer closers.Close(blobWriter)

	zipWriter := zip.NewWriter(blobWriter)
	defer closers.Close(zipWriter)

	return backup(ctx, zipWriter, conn)
}

func backup(ctx context.Context, writer *zip.Writer, conn *ftp.Conn) error {
	return conn.Walk(ctx, ftpPath, func(path string, info os.FileInfo, err error) error {
		switch {
		case err != nil:
			return err
		case info.IsDir():
			return nil
		}

		fileWriter, err := writer.CreateHeader(&zip.FileHeader{
			Name:     strings.TrimPrefix(path, "/"),
			Modified: info.ModTime(),
			Method:   zip.Deflate,
		})
		if err != nil {
			return err
		}

		fileReader, err := conn.NewReader(path)
		if err != nil {
			return err
		}
		defer closers.Close(fileReader)

		_, err = io.Copy(fileWriter, fileReader)
		return err
	})
}
