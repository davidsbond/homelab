package grafana

import (
	"context"
	"io"

	"github.com/grafana-tools/sdk"
)

type (
	// The Dashboard type is an io.Reader implementation that contains raw JSON representing
	// the dashboard.
	Dashboard struct {
		data []byte
		uid  string
		n    int
	}

	// The DashboardIterator is a function that is invoked for each dashboard when using Client.IterateDashboards.
	DashboardIterator func(ctx context.Context, d *Dashboard) error
)

// IterateDashboard searches for all dashboards in the grafana instance and invokes fn for each one. If fn returns
// an error or the context is cancelled iteration will stop.
func (cl *Client) IterateDashboards(ctx context.Context, fn DashboardIterator) error {
	boards, err := cl.grafana.Search(ctx, sdk.SearchType(sdk.SearchTypeDashboard))
	if err != nil {
		return err
	}

	for _, board := range boards {
		data, _, err := cl.grafana.GetRawDashboardByUID(ctx, board.UID)
		if err != nil {
			return err
		}

		db := &Dashboard{
			data: data,
			uid:  board.UID,
		}

		if err = fn(ctx, db); err != nil {
			return err
		}
	}

	return nil
}

// Read the content of the dashboard JSON, copying it to p.
func (d *Dashboard) Read(p []byte) (int, error) {
	if d.n >= len(d.data) {
		return 0, io.EOF
	}

	n := copy(p, d.data[d.n:])
	d.n += n
	return n, nil
}

// UID returns the unique UID for the dashboard.
func (d Dashboard) UID() string {
	return d.uid
}
