// Package filesystem contains methods for interacting with the filesystem.
package filesystem

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"os"
	"path/filepath"

	"pkg.dsb.dev/closers"
)

// Archive the desired directory, writing the contents as a tar.gz to the provided writer.
func Archive(ctx context.Context, dir string, w io.Writer) error {
	gzipWriter := gzip.NewWriter(w)
	tarWriter := tar.NewWriter(gzipWriter)

	defer closers.Close(gzipWriter)
	defer closers.Close(tarWriter)

	const mode = 0o600
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		if info.IsDir() {
			return nil
		}

		header := &tar.Header{
			Name: path,
			Mode: mode,
			Size: info.Size(),
		}

		if err = tarWriter.WriteHeader(header); err != nil {
			return err
		}

		file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return err
		}

		defer closers.Close(file)
		_, err = io.Copy(tarWriter, file)
		return err
	})
}
