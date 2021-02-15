package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"pkg.dsb.dev/app"
	"pkg.dsb.dev/server"

	"github.com/davidsbond/homelab/internal/health"
	"github.com/davidsbond/homelab/internal/health/scrape"
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(),
	)

	if err := a.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	r := mux.NewRouter()
	scraper := scrape.NewScraper(client)
	health.NewHTTP(scraper).Register(r)
	return server.ServeHTTP(ctx, server.WithHandler(r))
}
