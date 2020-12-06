package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gorilla/mux"
	"pkg.dsb.dev/app"
	"pkg.dsb.dev/flag"
	"pkg.dsb.dev/server"

	"github.com/davidsbond/homelab/internal/directory"
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(
			&flag.String{
				Name:        "config-path",
				Usage:       "Location of the configuration file",
				EnvVar:      "CONFIG_PATH",
				Value:       "./config.yaml",
				Destination: &configPath,
			},
		),
	)

	if err := a.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var configPath string

func run(ctx context.Context) error {
	config, err := directory.LoadConfig(configPath)
	if err != nil {
		return err
	}

	r := mux.NewRouter()
	directory.NewHTTP(config).Register(r)
	return server.ServeHTTP(ctx, server.WithHandler(r))
}
