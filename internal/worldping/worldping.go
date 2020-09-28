// Package worldping is used to perform a ping test for world-wide servers. It uses the debian
// FTP server mirrors to perform ping tests.
package worldping

import (
	"context"
	"sync"
	"time"

	"github.com/go-ping/ping"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/sync/errgroup"
	"pkg.dsb.dev/tracing"
)

// Run a test for all available world ping servers. Returns a map of the average
// RTT keyed by the lat-long value of the location.
func Run(ctx context.Context, privileged bool) (map[Server]time.Duration, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "world-ping")
	defer span.Finish()

	out := make(map[Server]time.Duration)
	mux := &sync.Mutex{}

	grp, ctx := errgroup.WithContext(ctx)
	for _, svr := range servers {
		performPing(ctx, grp, mux, svr, privileged, out)
	}

	return out, grp.Wait()
}

func performPing(ctx context.Context, grp *errgroup.Group, mux sync.Locker, server Server, privileged bool, out map[Server]time.Duration) {
	grp.Go(func() error {
		span, _ := opentracing.StartSpanFromContext(ctx, "ping-test")
		defer span.Finish()

		url := server.URL()

		span.SetTag("ping.host", url)
		span.SetTag("ping.name", server.Name)

		pinger, err := ping.NewPinger(url)
		if err != nil {
			return tracing.WithError(span, err)
		}

		pinger.Count = 3
		pinger.Timeout = time.Minute
		pinger.SetPrivileged(privileged)

		if err = pinger.Run(); err != nil {
			return tracing.WithError(span, err)
		}

		stats := pinger.Statistics()
		span.SetTag("ping.average_rtt", stats.AvgRtt)

		mux.Lock()
		out[server] = stats.AvgRtt
		mux.Unlock()

		return nil
	})
}
