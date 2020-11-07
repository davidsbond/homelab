package grafana

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type (
	// The DataSource type is an io.Reader implementation that contains raw JSON representing
	// the data source.
	DataSource struct {
		id   uint
		data []byte
		n    int
	}

	// The DataSourceIterator is a function that is invoked for each dashboard when using Client.IterateDataSources.
	DataSourceIterator func(ctx context.Context, d *DataSource) error
)

// IterateDataSources searches for all data sources in the grafana instance and invokes fn for each one. If fn returns
// an error or the context is cancelled iteration will stop.
func (cl *Client) IterateDataSources(ctx context.Context, fn DataSourceIterator) error {
	sources, err := cl.grafana.GetAllDatasources(ctx)
	if err != nil {
		return err
	}

	for _, source := range sources {
		data, err := json.Marshal(source)
		if err != nil {
			return err
		}

		ds := &DataSource{
			data: data,
			id:   source.ID,
		}

		if err = fn(ctx, ds); err != nil {
			return err
		}
	}

	return nil
}

// Read the content of the data source JSON, copying it to p.
func (d *DataSource) Read(p []byte) (int, error) {
	if d.n >= len(d.data) {
		return 0, io.EOF
	}

	n := copy(p, d.data[d.n:])
	d.n += n
	return n, nil
}

// ID returns the unique ID for the data source.
func (d DataSource) ID() string {
	return fmt.Sprint(d.id)
}
