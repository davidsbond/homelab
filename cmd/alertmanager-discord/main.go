package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gorilla/mux"
	"pkg.dsb.dev/app"
	"pkg.dsb.dev/flag"
	"pkg.dsb.dev/server"

	"github.com/davidsbond/homelab/internal/alertmanager"
	"github.com/davidsbond/homelab/internal/discord"
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(
			&flag.String{
				Name:        "discord-webhook-url",
				Usage:       "Discord webhook URL for forwarding alerts",
				EnvVar:      "DISCORD_WEBHOOK_URL",
				Required:    true,
				Destination: &discordWebhookURL,
			},
		),
	)

	if err := a.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var discordWebhookURL string

func run(ctx context.Context) error {
	r := mux.NewRouter()
	dispatcher := discord.NewAlertDispatcher(discordWebhookURL)
	alertmanager.NewWebhookHandler(dispatcher).Register(r)
	return server.ServeHTTP(ctx, server.WithHandler(r))
}
