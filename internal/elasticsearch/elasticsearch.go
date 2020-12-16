// Package elasticsearch contains utilities for interacting with elasticsearch clusters.
package elasticsearch

import (
	"context"
	"errors"
	"fmt"
	"time"

	"pkg.dsb.dev/health"

	"github.com/olivere/elastic/v7"
)

type (
	// The IndexCleaner type is used to delete indexes that fall under certain criteria.
	IndexCleaner struct {
		client *elastic.Client
	}
)

// NewIndexCleaner returns a new instance of the IndexCleaner type that communicates with the elasticsearch
// instance configured via the provided client.
func NewIndexCleaner(client *elastic.Client) *IndexCleaner {
	ic := &IndexCleaner{client: client}
	health.AddCheck("elasticsearch", ic.Ping)
	return ic
}

// ErrUnhealthy is the error given when Ping determines the client is un an unhealthy state.
var ErrUnhealthy = errors.New("elasticsearch status unhealthy")

// Ping returns a non-nil error if the cluster is deemed unhealthy.
func (ic *IndexCleaner) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	resp, err := ic.client.ClusterHealth().Do(ctx)
	switch {
	case err != nil:
		return err
	case resp.Status == "red":
		return fmt.Errorf("%w: %s", ErrUnhealthy, resp.Status)
	default:
		return nil
	}
}

// Clean indexes older than maxAge, whose names match those provided in formats. Each string in the formats slice
// should follow the standard go formatting directives. Each is generated using the year, month and day of today
// minus maxAge. Once indexes have been deleted, a force merge is performed on deleted indexes only.
func (ic *IndexCleaner) Clean(ctx context.Context, formats []string, maxAge time.Duration) error {
	date := time.Now().Add(-maxAge)

	indices := make([]string, len(formats))
	for i, format := range formats {
		indices[i] = fmt.Sprintf(format, date.Year(), int(date.Month()), date.Day())
	}

	var expunge bool
	for _, index := range indices {
		_, err := ic.client.DeleteIndex(index).Do(ctx)
		switch {
		case elastic.IsNotFound(err):
			continue
		case err != nil:
			return err
		default:
			expunge = true
		}
	}

	// Don't perform the force merge if nothing has been deleted.
	if !expunge {
		return nil
	}

	_, err := ic.client.Forcemerge().OnlyExpungeDeletes(true).Do(ctx)
	return err
}
