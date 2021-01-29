package main

import (
	"context"
	"fmt"
	"os"

	"github.com/davidsbond/homelab/internal/minecraft"

	"pkg.dsb.dev/app"
	"pkg.dsb.dev/flag"
	"pkg.dsb.dev/metrics"
)

func main() {
	a := app.New(
		app.WithRunner(run),
		app.WithFlags(
			&flag.String{
				Name:        "minecraft-ip",
				Usage:       "The IP address of the minecraft server",
				EnvVar:      "MINECRAFT_IP",
				Required:    true,
				Destination: &minecraftIP,
			},
		),
	)

	if err := a.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var minecraftIP string

func run(ctx context.Context) error {
	metrics.Register(
		playersOnline,
		maxPlayers,
	)

	client := minecraft.NewClient(minecraftIP)
	status, err := client.GetStatus(ctx)
	if err != nil {
		return err
	}

	playersOnline.Set(float64(status.Players.Online))
	maxPlayers.Set(float64(status.Players.Max))
	return metrics.Push()
}
