// Package ftp contains an FTP client implementation that supports health checks and tracing.
package ftp

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/opentracing/opentracing-go"

	"pkg.dsb.dev/health"
	"pkg.dsb.dev/multierror"
	"pkg.dsb.dev/tracing"
)

type (
	// The Conn type represents a connection to an FTP server.
	Conn struct {
		inner *ftp.ServerConn
	}

	// FileInfo is an os.FileInfo implementation that represents FTP file information.
	FileInfo struct {
		name    string
		size    int64
		isDir   bool
		modTime time.Time
	}

	config struct {
		username string
		password string
	}
)

// Open a connection to the specified FTP server, applying any provided options. See the Option type for configuration
// options.
func Open(ctx context.Context, addr string, opts ...Option) (*Conn, error) {
	c := &config{}
	for _, opt := range opts {
		opt(c)
	}

	conn, err := ftp.Dial(addr,
		ftp.DialWithContext(ctx),
		ftp.DialWithTimeout(time.Minute),
	)
	if err != nil {
		return nil, err
	}

	out := &Conn{inner: conn}
	health.AddCheck(addr, out.Ping)

	if c.username == "" {
		return out, out.Ping()
	}

	if err = conn.Login(c.username, c.password); err != nil {
		return nil, multierror.Append(err, out.Close())
	}

	return out, out.Ping()
}

// Close the connection to the FTP server.
func (c *Conn) Close() error {
	return c.inner.Quit()
}

// NewReader returns a new io.ReadCloser implementation that reads the content of the file at the specified path.
func (c *Conn) NewReader(path string) (io.ReadCloser, error) {
	return c.inner.Retr(path)
}

// ListDir returns a slice of os.FileInfo implementations for all files/directories found at the given path of the
// FTP server.
func (c *Conn) ListDir(ctx context.Context, path string) ([]os.FileInfo, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ftp-list-dir")
	span.SetTag("path", path)
	defer span.Finish()

	entries, err := c.inner.List(path)
	if err != nil {
		return nil, tracing.WithError(span, err)
	}

	span.SetTag("count", len(entries))
	out := make([]os.FileInfo, len(entries))
	for i, entry := range entries {
		out[i] = &FileInfo{
			name:    entry.Name,
			size:    int64(entry.Size),
			isDir:   entry.Type == ftp.EntryTypeFolder,
			modTime: entry.Time,
		}
	}

	return out, nil
}

// Walk walks the file tree rooted at root, calling fn for each file or
// directory in the tree. All errors that arise visiting files
// and directories are filtered by fn.
func (c *Conn) Walk(ctx context.Context, path string, fn filepath.WalkFunc) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ftp-walk")
	defer span.Finish()
	span.SetTag("path", path)

	walker := c.inner.Walk(path)
	for walker.Next() {
		entry := walker.Stat()
		err := fn(walker.Path(), &FileInfo{
			name:    entry.Name,
			size:    int64(entry.Size),
			isDir:   entry.Type == ftp.EntryTypeFolder,
			modTime: entry.Time,
		}, walker.Err())

		switch {
		case errors.Is(err, filepath.SkipDir):
			walker.SkipDir()
		case err != nil:
			return err
		}
	}

	return walker.Err()
}

// Ping asserts that the connection to the FTP server is alive and healthy.
func (c *Conn) Ping() error {
	return c.inner.NoOp()
}

// Name returns the file name.
func (f *FileInfo) Name() string {
	return f.name
}

// Size returns the file size in bytes.
func (f *FileInfo) Size() int64 {
	return f.size
}

// Mode returns os.ModeDir if the os.FileMode implementation is a directory. Otherwise it returns
// os.ModeIrregular.
func (f *FileInfo) Mode() os.FileMode {
	if f.isDir {
		return os.ModeDir
	}

	return os.ModeIrregular
}

// ModTime returns the last modified time of the file.
func (f *FileInfo) ModTime() time.Time {
	return f.modTime
}

// IsDir returns true if the current file is a directory.
func (f *FileInfo) IsDir() bool {
	return f.isDir
}

// Sys returns nil.
func (f *FileInfo) Sys() interface{} {
	return nil
}
