// Package app contains methods for configuring a command-line application.
package app

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	"pkg.dsb.dev/closers"
	"pkg.dsb.dev/environment"
	"pkg.dsb.dev/health"
	"pkg.dsb.dev/logging"
	"pkg.dsb.dev/metrics"
	"pkg.dsb.dev/monitoring"
	"pkg.dsb.dev/tracing"
)

type (
	// The App type represents a command-line application.
	App struct {
		inner *cli.App
	}

	// The Option type is a function that can apply configuration to the command-line
	// application.
	Option func(app *cli.App)

	// The RunFunc type describes a method invoked to start a cli command.
	RunFunc func(ctx context.Context) error
)

// New creates a new command-line application with no action. Metadata such as name, version etc. are set by default.
func New(opts ...Option) *App {
	app := cli.NewApp()
	app.Version = environment.Version
	app.Name = environment.ApplicationName
	app.Usage = environment.ApplicationDescription
	app.Compiled = environment.Compiled()
	app.Commands = []*cli.Command{
		{
			Name:  "check",
			Usage: "Performs an HTTP health check",
			Action: func(c *cli.Context) error {
				return health.PerformCheck(c.Context)
			},
		},
	}

	for _, opt := range opts {
		opt(app)
	}

	return &App{inner: app}
}

// Run the application.
func (a *App) Run() error {
	ctx, cancel := context.WithCancel(environment.NewContext())
	defer cancel()

	action := a.inner.Action
	a.inner.Action = func(c *cli.Context) error {
		logging.Init()

		svr := operationalServer()
		var grp errgroup.Group

		// Setup metrics/health servers
		grp.Go(svr.ListenAndServe)
		grp.Go(func() error {
			<-ctx.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			return svr.Shutdown(ctx)
		})

		// Setup tracing
		closer, err := tracing.New()
		if err != nil {
			return err
		}
		defer closers.Close(closer)

		// Setup monitoring
		monitor, err := monitoring.New()
		if err != nil {
			return err
		}
		defer closers.Close(monitor)

		return action(c)
	}

	err := a.inner.RunContext(ctx, os.Args)
	return monitoring.WithError(err)
}

// WithRunner sets the action function to be used by the command-line application when Run is
// called.
func WithRunner(run RunFunc) Option {
	return func(app *cli.App) {
		app.Action = func(c *cli.Context) error {
			return run(c.Context)
		}
	}
}

// WithFlags sets command-line flags that can be applied before Run is called.
func WithFlags(flags ...cli.Flag) Option {
	return func(app *cli.App) {
		app.Flags = flags
		app.Flags = append(app.Flags, tracing.Flags...)
		app.Flags = append(app.Flags, monitoring.Flags...)
		app.Flags = append(app.Flags, health.Flags...)
		app.Flags = append(app.Flags, metrics.Flags...)
		app.Flags = append(app.Flags, logging.Flags...)

		sort.Sort(cli.FlagsByName(app.Flags))
	}
}

func operationalServer() *http.Server {
	r := mux.NewRouter()
	health.ServeHTTP(r)
	metrics.Serve(r)

	return &http.Server{
		Addr:     ":8081",
		Handler:  r,
		ErrorLog: log.New(ioutil.Discard, "", log.Flags()),
	}
}
