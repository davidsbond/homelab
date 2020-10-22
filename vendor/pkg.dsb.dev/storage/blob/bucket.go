// Package blob is used to read/write data from blob stores.
package blob

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"github.com/opentracing/opentracing-go"
	"gocloud.dev/blob"
	"gocloud.dev/gcerrors"

	"pkg.dsb.dev/health"
	"pkg.dsb.dev/tracing"

	// Gocloud driver for in-memory buckets.
	_ "gocloud.dev/blob/memblob"
	// Gocloud driver for s3 buckets.
	_ "gocloud.dev/blob/s3blob"
	// Gocloud driver for file buckets.
	_ "gocloud.dev/blob/fileblob"
)

// ErrNotExist is the error returned when attempting to query blob data that
// does not exist.
var ErrNotExist = errors.New("not exist")

type (
	// The Bucket type represents a bucket of blob data that can be arbitrarily
	// written to and read using a key.
	Bucket struct {
		inner *blob.Bucket
	}

	// The Writer type is able to write blob data for a desired key. It is instrumented
	// to write spans to the tracer when blob data is written.
	Writer struct {
		ctx   context.Context
		inner *blob.Writer
		key   string
	}

	// The Reader type is able to read blob data for a desired key. It is instrumented
	// to write spans to the tracer when blob data is read.
	Reader struct {
		ctx   context.Context
		span  opentracing.Span
		inner *blob.Reader
		n     int
		key   string
	}

	// The Blob type contains metadata on an item in the blob store.
	Blob struct {
		Key     string
		Size    int64
		ModTime time.Time
	}

	// Iterator is a function used on a call to Bucket.Iterate that is invoked for
	// each item in the bucket.
	Iterator func(ctx context.Context, item Blob) error
)

// OpenBucket opens the bucket identified by the URL given.
func OpenBucket(ctx context.Context, dsn string) (*Bucket, error) {
	dsn = os.ExpandEnv(dsn)

	inner, err := blob.OpenBucket(ctx, dsn)
	if err != nil {
		return nil, err
	}

	bkt := &Bucket{inner: inner}
	health.AddCheck(dsn, bkt.Ping)
	return bkt, bkt.Ping()
}

// Close the connection to the bucket.
func (bkt *Bucket) Close() error {
	return bkt.inner.Close()
}

// Ping is a basic check to ensure the connection to the bucket
// is healthy.
func (bkt *Bucket) Ping() error {
	it := bkt.inner.List(&blob.ListOptions{})
	_, err := it.Next(context.Background())
	switch {
	case errors.Is(err, io.EOF):
		return nil
	default:
		return err
	}
}

// ErrStopIterating is the error used to stop the iterator from continuing.
var ErrStopIterating = errors.New("stop iterating")

// Iterate over the contents of the bucket, invoking fn for each item, excluding any directories.
// Iteration can be cancelled via the provided context or by fn returning ErrStopIterating or any other
// non-nil error.
func (bkt *Bucket) Iterate(ctx context.Context, fn Iterator) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blob-iterate")
	defer span.Finish()

	iterator := bkt.inner.List(nil)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			item, err := iterator.Next(ctx)
			switch {
			case errors.Is(err, io.EOF):
				return nil
			case err != nil:
				return tracing.WithError(span, err)
			}

			if item.IsDir {
				continue
			}

			bl := Blob{
				Key:     item.Key,
				Size:    item.Size,
				ModTime: item.ModTime,
			}

			err = fn(ctx, bl)
			switch {
			case errors.Is(err, ErrStopIterating):
				return nil
			case err != nil:
				return tracing.WithError(span, err)
			}
		}
	}
}

// Delete deletes the blob stored at key. Returns ErrNotExist if the key does not
// exist.
func (bkt *Bucket) Delete(ctx context.Context, key string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blob-delete")
	span.SetTag("blob.key", key)
	defer span.Finish()

	err := bkt.inner.Delete(ctx, key)
	switch {
	case gcerrors.Code(err) == gcerrors.NotFound:
		return tracing.WithError(span, ErrNotExist)
	case err != nil:
		return tracing.WithError(span, err)
	default:
		return nil
	}
}

// NewWriter returns a Writer that writes to the blob stored at key.
func (bkt *Bucket) NewWriter(ctx context.Context, key string) (*Writer, error) {
	inner, err := bkt.inner.NewWriter(ctx, key, nil)
	if err != nil {
		return nil, err
	}

	blobsOpen.Inc()
	return &Writer{inner: inner, ctx: ctx, key: key}, nil
}

// NewRangeReader returns a Reader to read content from the blob stored at key from one offset to a certain
// number of bytes.
func (bkt *Bucket) NewRangeReader(ctx context.Context, key string, offset, length int64) (*Reader, error) {
	inner, err := bkt.inner.NewRangeReader(ctx, key, offset, length, nil)
	switch {
	case gcerrors.Code(err) == gcerrors.NotFound:
		return nil, ErrNotExist
	case err != nil:
		return nil, err
	default:
		blobsOpen.Inc()
		return &Reader{inner: inner, ctx: ctx, key: key}, nil
	}
}

// NewReader returns a Reader to read content from the blob stored at key.
func (bkt *Bucket) NewReader(ctx context.Context, key string) (*Reader, error) {
	inner, err := bkt.inner.NewReader(ctx, key, nil)
	switch {
	case gcerrors.Code(err) == gcerrors.NotFound:
		return nil, ErrNotExist
	case err != nil:
		return nil, err
	default:
		blobsOpen.Inc()
		return &Reader{inner: inner, ctx: ctx, key: key}, nil
	}
}

// Close the blob reader, this will finish the traced span.
func (r *Reader) Close() error {
	bytesRead.Add(float64(r.n))
	r.span.SetTag("blob.bytes_read", r.n)
	r.span.Finish()
	blobsOpen.Dec()
	return r.inner.Close()
}

// Close the blob writer.
func (w *Writer) Close() error {
	blobsOpen.Dec()
	return w.inner.Close()
}

// Write the bytes in p to the blob store.
func (w *Writer) Write(p []byte) (int, error) {
	span, _ := opentracing.StartSpanFromContext(w.ctx, "blob-write")
	span.SetTag("blob.key", w.key)
	defer span.Finish()

	n, err := w.inner.Write(p)
	switch {
	case err != nil:
		return n, tracing.WithError(span, err)
	default:
		bytesWritten.Add(float64(n))
		span.SetTag("blob.bytes_written", n)
		return n, nil
	}
}

// ReadFrom reads from r and writes to w until EOF or error.
// The return value is the number of bytes read from r.
func (w *Writer) ReadFrom(r io.Reader) (int64, error) {
	span, _ := opentracing.StartSpanFromContext(w.ctx, "blob-write")
	span.SetTag("blob.key", w.key)
	defer span.Finish()

	n, err := w.inner.ReadFrom(r)
	switch {
	case err != nil:
		return n, tracing.WithError(span, err)
	default:
		bytesWritten.Add(float64(n))
		span.SetTag("blob.bytes_written", n)
		return n, nil
	}
}

// Read implements io.Reader (https://golang.org/pkg/io/#Reader).
func (r *Reader) Read(p []byte) (int, error) {
	r.span, _ = opentracing.StartSpanFromContext(r.ctx, "blob-read")
	r.span.SetTag("blob.key", r.key)

	n, err := r.inner.Read(p)
	switch {
	case errors.Is(err, io.EOF):
		return n, err
	case err != nil:
		return n, tracing.WithError(r.span, err)
	default:
		r.n += n
		return n, nil
	}
}
